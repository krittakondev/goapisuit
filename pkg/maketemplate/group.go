package maketemplate

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/krittakondev/goapisuit/pkg/utils"
)

func NewGroup(name, path string) (createPathRoute string, err error) {
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
	if err = tmplRoute.Execute(file, NewMakeTemplate(name)); err != nil {
		return 
	}
	filetmp, err := os.OpenFile(".tmpgroups", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open file: %v", err)
	}
	defer filetmp.Close()
	_, err = filetmp.WriteString(strings.TrimPrefix(strings.ReplaceAll(path, "/init_suit.go", "") + "\n", "internal/routes"))

	return
}
func (mr *GroupsLoader) NewGroupLoader() (err error) {
	tmplGroupLoader, err := template.New("GroupsLoader").Parse(templateGroupsSetup)
	if err != nil{
		return
	}
	file, err := os.OpenFile("internal/setup/groupsloader.go", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 
	}
	defer file.Close()
	if err = tmplGroupLoader.Execute(file, mr); err != nil {
		return 
	}

	return
}

func migrateGroup(){
	group, err := os.ReadFile(".tmpgroups")
	if err != nil{
		log.Fatal(err)
	}
	list := strings.Split(string(group), "\n")

	project_path, _ := utils.GetProjectName()
	list_call := CreateTemplateGroupsSetupCall(list)
	list_import := CreateTemplateGroupsSetupImport(project_path, list)
	mktemp := &GroupsLoader{
		ImportRouteGroup: strings.Join(list_import, "\n"),
		SetupGroups:  strings.Join(list_call, "\n"),
	}
	err = mktemp.NewGroupLoader()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateInitSuitInGroup(group_name string) {
	re := regexp.MustCompile("/+")
	group_name = re.ReplaceAllString(group_name, "/")
	list_group := strings.Split(group_name, "/")
	info, _ := os.Stat("internal/routes")
	if !info.IsDir() {
		log.Println("Don't have internal/routes path please check your current path")
	}
	err := os.MkdirAll(re.ReplaceAllString("internal/routes/"+group_name, "/"), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}
	path := "internal/routes"
	for _, val := range list_group {
		path += "/"+val
		path_init := re.ReplaceAllString(path + "/init_suit.go", "/")
		if _, err := os.Stat(path_init); err != nil {
			if os.IsNotExist(err) {
				if created, err1 := NewGroup(utils.KebabToCamel(val), path_init); err1 != nil {
					log.Fatal(err1)
				} else {
					fmt.Printf("created %s\n", created)

				}

			} else {
				log.Println(path_init + "exist!")
			}

		}
	}
	migrateGroup()
}
