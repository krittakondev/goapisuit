package maketemplate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/krittakondev/goapisuit/pkg/utils"
)

type MakeRoute struct {
	Name string
}

type migrate struct {
	Case        string
	ProjectName string
}

func recreateMigrateDbFunc() error {
	read, err := os.ReadFile(".tmpmodels")
	if err != nil {
		return err
	}
	var tm migrate
	tm.ProjectName, err = utils.GetProjectName()
	if err != nil {
		return err
	}
	split_name := strings.Split(string(read), "\n")
	for _, val := range split_name {
		if val != "" {
			tm.Case += fmt.Sprintf("case strings.ToLower(\"%s\"):\n", val)
			tm.Case += fmt.Sprintf("\treturn db.AutoMigrate(&models.%s{})\n", val)
		}
	}

	tmpl, err := template.ParseFiles(filepath.Join("pkg","maketemplate","tmpl","migrate.go.tmpl"))
	createPath := "internal/database/migrate.go"
	if err != nil {
		return err
	}
	file, err := os.OpenFile(createPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(file, tm); err != nil {
		return err
	}
	return nil
}

func (mr *MakeRoute) New() (arrPath []string, err error) {
	tmplRoute, err := template.ParseFiles(filepath.Join("pkg","maketemplate","tmpl","route.go.tmpl"))
	if err != nil {
		return 
	}
	tmplModel, err := template.ParseFiles(filepath.Join("pkg","maketemplate","tmpl","model.go.tmpl"))
	if err != nil {
		return 
	}
	createPathRoute := "internal/routes/" + mr.Name + ".go"
	createPathModel := "internal/models/" + mr.Name + ".go"

	if _, err = os.Stat(createPathRoute); err == nil {
		// log.Fatal(createPathRoute + " is Exist")
		return
	}
	if _, err = os.Stat(createPathModel); err == nil {
		// log.Fatal(createPathModel + " is Exist")
		return
	}

	file, err := os.OpenFile(createPathRoute, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 
	}
	if err = tmplRoute.Execute(file, mr); err != nil {
		return 
	}

	file, err = os.OpenFile(createPathModel, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if err = tmplModel.Execute(file, mr); err != nil {
		return 
	}
	file, err = os.OpenFile(".tmpmodels", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open file: %v", err)
	}
	defer file.Close()

	if _, err = file.WriteString(mr.Name + "\n"); err != nil {
		return 
	}

	if err = recreateMigrateDbFunc(); err != nil {
		return
	}

	arrPath = []string{
		createPathModel,
		createPathRoute,
	}
	return 
}
