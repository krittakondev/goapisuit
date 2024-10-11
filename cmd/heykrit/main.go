package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
	"github.com/krittakondev/goapisuit"
	"github.com/krittakondev/goapisuit/database"
	"github.com/krittakondev/goapisuit/pkg/maketemplate"
	"github.com/krittakondev/goapisuit/pkg/utils"
	"github.com/manifoldco/promptui"
)

func argsScan(opt_name string, args ...*string) (err error) {
	len_scan := 2 + len(args)
	if len(os.Args) < len_scan {
		fmt.Printf(`Error: Missing %d required arguments.
Usage: heykrit %s`, len_scan-len(os.Args), opt_name)
		fmt.Println()
		os.Exit(1)
	}
	for i, _ := range args {
		*args[i] = os.Args[2+i]
	}
	return
}

func createInitSuitInGroup(arr []string) {
	path := "internal/routes"
	for _, val := range arr {
		path += "/"+val
		path_init := path + "/init_suit.go"
		if _, err := os.Stat(path_init); err != nil {
			if os.IsNotExist(err) {
				PathProject, _ := utils.GetProjectName()
				mkroute := &maketemplate.MakeRoute{
					Name:        utils.KebabToCamel(val),
					PathProject: PathProject,
				}
				if created, err1 := mkroute.NewGroup(path_init); err1 != nil {
					log.Fatal(err1)
				} else {
					fmt.Printf("created %s\n", created)

				}

			} else {
				log.Println(path_init + "exist!")
			}

		}
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass a command")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		if _, err := os.Stat("go.mod"); err == nil {
			projectPath, err := utils.GetProjectName()
			if err != nil {
				log.Fatal(err)
			}

			promptAskCreateProject := promptui.Select{
				Label: fmt.Sprintf("Do you want init goapisuit for %s?", projectPath),
				Items: []string{"NO", "YES"},
			}
			_, result, err := promptAskCreateProject.Run()
			if err != nil {
				log.Fatal(err)
			}
			if result != "YES" {
				fmt.Println("exit init")
				os.Exit(0)
			}
			useDocker := promptui.Select{
				Label: fmt.Sprintf("Do you want use Docker?"),
				Items: []string{"NO", "YES"},
			}
			_, resultuseDocker, err := useDocker.Run()
			if err != nil {
				log.Fatal(err)
			}
			template := &maketemplate.Template{
				ProjectName: projectPath,
			}
			template.EnvStruct.AppPort = "3000"
			template.EnvStruct.DbDatabase = "goapisuit"
			template.EnvStruct.DbUsername = "suit"
			template.EnvStruct.DbHost = "127.0.0.1"
			template.EnvStruct.DbPort = "3306"

			if resultuseDocker == "YES" {
				dbpassword := os.Getenv("DB_PASSWORD")
				if len(dbpassword) > 0 {
					template.EnvStruct.DbPassword = dbpassword
				} else {
					template.EnvStruct.DbPassword, _ = utils.GenerateSecret(32)
				}
			}

			done := make(chan bool)
			go template.InitProject(done, resultuseDocker == "YES")
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Start()
			<-done
			fmt.Println("done")
			s.Stop()

		} else if os.IsNotExist(err) {
			fmt.Println("Not found go.mod")
		} else {
			log.Fatalf("Error checking go.mod: %v\n", err)
		}

	case "make":
		var routeName string
		argsScan("make [...args]", &routeName)
		PathProject, _ := utils.GetProjectName()
		mkroute := &maketemplate.MakeRoute{
			Name:        utils.KebabToCamel(routeName),
			PathProject: PathProject,
		}
		if arr, err := mkroute.New(); err != nil {
			log.Fatal(err)
		} else {
			for _, str := range arr {
				fmt.Printf("created %s\n", str)
			}

		}
	case "make:group":
		var group_name string
		argsScan(command+" [...args]", &group_name)
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
		createInitSuitInGroup(list_group)

	case "make:route":
		var routeName string
		argsScan(command+" [...args]", &routeName)
		PathProject, _ := utils.GetProjectName()
		mkroute := &maketemplate.MakeRoute{
			Name:        utils.KebabToCamel(routeName),
			PathProject: PathProject,
		}
		if str, err := mkroute.NewRoute(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("created %s\n", str)
		}
	case "make:model":
		var routeName string
		argsScan(command+" [...args]", &routeName)
		PathProject, _ := utils.GetProjectName()
		mkroute := &maketemplate.MakeRoute{
			Name:        utils.KebabToCamel(routeName),
			PathProject: PathProject,
		}
		if str, err := mkroute.NewModel(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("created %s\n", str)
		}
	case "db:testconnect":
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
		if _, err := database.MysqlConnect(); err != nil {
			log.Fatal(err)
		}
		log.Print("connect success")
	case "db:migrate":
		var model_name string
		argsScan(command+" [...args]", &model_name)

		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
		// _, err := database.MysqlConnect()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		tmpmodels, _ := goapisuit.LoadTmpModel()
		var ModelName string
		for _, str := range tmpmodels {
			if strings.ToLower(model_name) == strings.ToLower(str) {
				ModelName = str
				break
			}
		}
		if ModelName == "" {
			log.Fatalf("Not found model %s\n", model_name)
		}

		fmt.Printf("Do you want migrate %s Model? [y/N]:", model_name)
		Ans := "n"
		fmt.Scanf("%s\n", &Ans)
		if strings.ToLower(Ans) != "y" {
			log.Fatal("not migrate!")
		}
		pathProject, err := utils.GetProjectName()
		if err != nil {
			log.Fatal(pathProject)
		}
		mg := maketemplate.Migrate{
			PathProject: pathProject,
			Name:        ModelName,
		}

		if err := mg.Migrate(); err != nil {
			log.Fatal(err)
		}

		// err = database.Migrate(db, model_name)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		log.Printf("%s Migrate Success", model_name)

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
