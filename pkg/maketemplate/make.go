package maketemplate

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/krittakondev/goapisuit/pkg/utils"
)

type MakeTemplate struct {
	Name string
	PathProject string
	ModelName string
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

func NewMakeTemplate(name string) *MakeTemplate{
	project_name, _ := utils.GetProjectName()
	
	return &MakeTemplate{
		PathProject: project_name,
		Name: name,
	}
}

func NewModel(name string) (createPathModel string, err error) {
	tmplModel, err  := template.New("model").Parse(templateMakeModel)
	if err != nil{
		return
	}
	model_name := utils.PathToCamelCase(name)
	model_name = strings.Join(strings.Split(model_name, "-"), "_")
	re := regexp.MustCompile("/+")
	createPathModel = re.ReplaceAllString("internal/models/" + model_name + ".go", "/")
	if _, err = os.Stat(createPathModel); err == nil {
		err = errors.New(createPathModel + " is Exist")
		return 
	}
	file, err := os.OpenFile(createPathModel, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	if err = tmplModel.Execute(file, NewMakeTemplate(model_name)); err != nil {
		return 
	}
	filetmp, err := os.OpenFile(".tmpmodels", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open file: %v", err)
	}
	defer filetmp.Close()
	_, err = filetmp.WriteString(model_name + "\n")
	return
}

func NewRoute(name string) (createPathRoute string, err error) {
	tmplRoute, err := template.New("route").Parse(templateMakeRouter)
	if err != nil{
		return
	}
	re := regexp.MustCompile("/+")
	createPathRoute = re.ReplaceAllString("internal/routes/" + name + ".go", "/")
	
	split_path_group := strings.Split(name, "/")
	path_group := strings.Join(split_path_group[:len(split_path_group)-1], "/")

	if path_group != ""{
		CreateInitSuitInGroup(path_group)
	}

	if _, err = os.Stat(createPathRoute); err == nil {
		err = errors.New(createPathRoute + " is Exist")
		return 
	}

	file, err := os.OpenFile(createPathRoute, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 
	}
	defer file.Close()
	name_route := utils.PathToCamelCase(split_path_group[len(split_path_group)-1])
	tmp :=  NewMakeTemplate(name_route)
	tmp.ModelName = utils.PathToModelFormatName(name)

	if err = tmplRoute.Execute(file, tmp); err != nil {
		return 
	}

	return
}
func New(name string) (arrPath []string, err error) {
	createPathModel, err :=  NewModel(name)
	if err != nil { 
		return
	}
	createPathRoute, err :=  NewRoute(name)
	if err != nil {
		return
	}

	arrPath = []string{
		createPathModel,
		createPathRoute,
	}
	return 
}
