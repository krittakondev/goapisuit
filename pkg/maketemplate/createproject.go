package maketemplate

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"text/template"

	"github.com/krittakondev/goapisuit"
	"github.com/krittakondev/goapisuit/pkg/utils"
)


type Template struct {
	EnvStruct
	GroupsLoader
	ProjectName string
}

func (t *Template) InitProject(done chan bool, useDocker bool) {
	arrDir := []string{
		"./cmd",
		"./public",
		"./internal/routes",
		"./internal/models",
		"./internal/database",
		"./internal/setup",
	}
	arrFile := map[string]string{
		".env":                           templateEnv,
		"./internal/routes/init_suit.go": templateRouter,
		"./internal/setup/groupsloader.go": templateGroupsSetup,
		"./cmd/server.go":                templateServer,
		"./public/index.html":            templatePublicIndex,
	}
	arrDockertemplate := map[string]string{
		"Dockerfile":         templateDockerfile,
		"docker-compose.yml": templateDockerCompose,
	}
	if useDocker {
		for key, value := range arrDockertemplate {
			arrFile[key] = value
		}

	}
	for _, dir := range arrDir {
		os.MkdirAll(dir, os.ModePerm)
	}

	t.EnvStruct.AppName = t.ProjectName

	t.EnvStruct.JwtSecret, _ = utils.GenerateSecret(32)

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
	fmt.Println("installing...")
	cmd := exec.Command("go", "get", "github.com/krittakondev/goapisuit@"+goapisuit.Version)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(out)
	}

	done <- true
}
