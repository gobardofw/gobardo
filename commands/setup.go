package commands

import (
	"fmt"
	"os"

	"github.com/gobardofw/crypto"
	"github.com/gobardofw/gobardo/internal/helpers"
	"github.com/gobardofw/gobardo/internal/questions"
	"github.com/google/uuid"
)

func setup(name string, w *questions.Wizard) {
	// helpers
	pathResolver := func(p string) string {
		return fmt.Sprintf("./%s/%s", name, p)
	}

	// init global data
	data := make(helpers.TemplateData)
	data["name"] = name
	data["description"] = w.Result("description")
	data["namespace"] = w.Result("namespace")
	data["locale"] = w.Result("locale")
	data["config"] = w.Result("config")
	data["cache"] = w.Result("cache")
	data["database"] = w.Result("database")
	data["translator"] = w.Result("translator")
	data["web"] = w.Result("web")

	// set app key
	c := crypto.NewCryptography(uuid.New().String())
	appKey, err := c.Hash(uuid.New().String(), crypto.SHA3256)
	helpers.Handle(err)
	data["appKey"] = appKey

	// Clean and compile
	os.RemoveAll(pathResolver("go.sum"))
	helpers.Handle(helpers.CompileTemplate(pathResolver("go.mod"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("app.tpl.go"), data))
	helpers.Handle(os.Rename(pathResolver("app.go"), pathResolver(fmt.Sprintf("%s.go", name))))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/app.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/boot.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/cache.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/config.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/crypto.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/logger.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/translator.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/validator.tpl.go"), data))
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/helpers/vars.tpl.go"), data))

	// config
	switch w.Result("config") {
	case "env":
		helpers.Handle(helpers.CompileTemplate(pathResolver("config/config.tpl.env"), data))
	case "json":
		helpers.Handle(helpers.CompileTemplate(pathResolver("config/config.tpl.json"), data))
	}
	helpers.Handle(helpers.CompileTemplate(pathResolver("src/config/config.tpl.go"), data))

	if w.Result("translator") == "memory" {
		os.RemoveAll(pathResolver("config/strings"))
	} else {
		os.Rename(pathResolver("config/strings/locale"), pathResolver("config/strings/")+w.Result("locale"))
	}

	if w.Result("config") == "memory" && w.Result("translator") == "memory" {
		os.RemoveAll(pathResolver("config"))
	}

	if w.Result("database") == "n" {
		os.RemoveAll(pathResolver("database"))
	} else {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/database.tpl.go"), data))

	}

	if w.Result("web") == "n" {
		os.RemoveAll(pathResolver("public"))
		os.RemoveAll(pathResolver("src/http"))
	} else {
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/commands/serve.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/bootstrap/web.tpl.go"), data))
		helpers.Handle(helpers.CompileTemplate(pathResolver("src/http/routes.tpl.go"), data))
	}
}
