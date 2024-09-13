package maketemplate

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"text/template"
)

type MakeRoute struct {
	Name string
	PathProject string
}

type Migrate struct {
	Name        string
	PathProject string
}

func (mg *Migrate) Migrate() error {

	// TODO: change to variable template
	tmpl, err := template.New("dbmigrate").Parse(templateDbMigrate)
	createPath := "tmp_migrate.go"
	if err != nil {
		return err
	}
	file, err := os.OpenFile(createPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(file, mg); err != nil {
		return err
	}
	cmd := exec.Command("go", "run", createPath)
	if out, err := cmd.CombinedOutput(); err != nil{
		return errors.New(string(out))
	}
	os.Remove(createPath)

	return nil
}

func (mr *MakeRoute) New() (arrPath []string, err error) {
	tmplRoute, err := template.New("route").Parse(templateMakeRouter)
	if err != nil{
		return
	}
	tmplModel, err  := template.New("model").Parse(templateMakeModel)
	if err != nil{
		return
	}

	createPathRoute := "internal/routes/" + mr.Name + ".go"
	createPathModel := "internal/models/" + mr.Name + ".go"

	if _, err = os.Stat(createPathRoute); err == nil {
		err = errors.New(createPathRoute + " is Exist")
		return 
	}
	if _, err = os.Stat(createPathModel); err == nil {
		err = errors.New(createPathModel + " is Exist")
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
		return
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

	arrPath = []string{
		createPathModel,
		createPathRoute,
	}
	return 
}
