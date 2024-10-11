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
	createPath := "tmp_migrate.go"
	tmpl, err := template.New("dbmigrate").Parse(templateDbMigrate)
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
		os.Remove(createPath)
		return errors.New(string(out))
	}
	os.Remove(createPath)

	return nil
}
func (mr *MakeRoute) NewModel() (createPathModel string, err error) {
	tmplModel, err  := template.New("model").Parse(templateMakeModel)
	if err != nil{
		return
	}
	createPathModel = "internal/models/" + mr.Name + ".go"
	if _, err = os.Stat(createPathModel); err == nil {
		err = errors.New(createPathModel + " is Exist")
		return 
	}
	file, err := os.OpenFile(createPathModel, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	if err = tmplModel.Execute(file, mr); err != nil {
		return 
	}
	filetmp, err := os.OpenFile(".tmpmodels", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open file: %v", err)
	}
	defer filetmp.Close()
	_, err = filetmp.WriteString(mr.Name + "\n")
	return
}

func (mr *MakeRoute) NewRoute() (createPathRoute string, err error) {
	tmplRoute, err := template.New("route").Parse(templateMakeRouter)
	if err != nil{
		return
	}
	createPathRoute = "internal/routes/" + mr.Name + ".go"

	if _, err = os.Stat(createPathRoute); err == nil {
		err = errors.New(createPathRoute + " is Exist")
		return 
	}

	file, err := os.OpenFile(createPathRoute, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 
	}
	defer file.Close()
	if err = tmplRoute.Execute(file, mr); err != nil {
		return 
	}

	return
}
func (mr *MakeRoute) NewGroup(path string) (createPathRoute string, err error) {
	tmplRoute, err := template.New("group").Parse(templateRouter)
	if err != nil{
		return
	}
	createPathRoute = path

	if _, err = os.Stat(createPathRoute); err == nil {
		err = errors.New(createPathRoute + " is Exist")
		return 
	}

	file, err := os.OpenFile(createPathRoute, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 
	}
	defer file.Close()
	if err = tmplRoute.Execute(file, mr); err != nil {
		return 
	}

	return
}
func (mr *MakeRoute) New() (arrPath []string, err error) {
	createPathModel, err :=  mr.NewModel()
	if err != nil {
		return
	}
	createPathRoute, err :=  mr.NewRoute()
	if err != nil {
		return
	}

	arrPath = []string{
		createPathModel,
		createPathRoute,
	}
	return 
}
