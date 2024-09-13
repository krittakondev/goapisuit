package maketemplate

import (
	"log"
	"os"

	"text/template"

	"github.com/krittakondev/goapisuit/pkg/utils"
)

type Template struct {
	envStruct
	ProjectName string
}

func (t *Template) InitProject(done chan bool) {
	arrDir := []string{
		"./cmd",
		"./public",
		"./internal/routes",
		"./internal/models",
		"./internal/database",
	}
	arrFile := map[string]string{
		"./internal/routes/init_suit.go": templateRouter,
		"./cmd/server.go":                templateServer,
		".env":                           templateEnv,
		"./public/index.html":            templatePublicIndex,
	}
	for _, dir := range arrDir {
		os.MkdirAll(dir, os.ModePerm)
	}

	t.envStruct.AppName = t.ProjectName
	t.envStruct.JwtSecret, _ = utils.GenerateSecret(32)
	for path, filedata := range arrFile {
		if _, err := os.Stat(path); err == nil {
			log.Printf("exist: %s\n", path)
		} else if os.IsNotExist(err) {
			tmpl, err := template.New(path).Parse(filedata)

			if err != nil {
				log.Printf("template error: %s\n", err)
			}
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0640)
			if err != nil {
				log.Println(err)
			}

			if err := tmpl.Execute(file, t); err != nil {
				log.Println(err)
			} else {
				log.Printf("created: %s\n", path)
			}
		}

	}

	done <- true
}
