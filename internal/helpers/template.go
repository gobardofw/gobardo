package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData structure
type TemplateData map[string]interface{}

// CompileTemplate compile template file
func CompileTemplate(filePath string, maps ...TemplateData) error {
	var data TemplateData
	if len(maps) == 0 {
		data = make(TemplateData)
	} else {
		data = maps[0]
		for i, m := range maps {
			if i == 0 {
				continue
			}
			for k, v := range m {
				data[k] = v
			}
		}
	}

	// Read file
	replacer := strings.NewReplacer("__anonymous__", data["namespace"].(string), "__goapp__", data["name"].(string))
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(bytes)
	content = replacer.Replace(content)
	dest := strings.ReplaceAll(filePath, ".tpl.", ".")

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	tpl, err := template.New(filepath.Base(filePath)).Delims(`// {{`, "}}").Parse(content)
	if err != nil {
		return err
	}

	if err := tpl.Execute(out, data); err != nil {
		return err
	}
	return nil
}
