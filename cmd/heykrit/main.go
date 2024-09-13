package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
	"github.com/krittakondev/goapisuit/internal/database"
	"github.com/krittakondev/goapisuit/pkg/maketemplate"
	"github.com/krittakondev/goapisuit/pkg/utils"
	"github.com/manifoldco/promptui"
)

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
			done := make(chan bool)
			if result == "YES" {
				template := &maketemplate.Template{
					ProjectName: projectPath,
				}
				go template.InitProject(done)
				s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
				s.Start()
				<- done
				fmt.Println("done")
				s.Stop()
			}

		} else if os.IsNotExist(err) {
			fmt.Println("Not found go.mod")
		} else {
			log.Fatalf("Error checking go.mod: %v\n", err)
		}

	case "make":
		if len(os.Args) < 3 {
			fmt.Println("Please Enter Route name")
			os.Exit(1)
		}
		routeName := os.Args[2]
		PathProject, _ := utils.GetProjectName()
		mkroute := &maketemplate.MakeRoute{
			Name: utils.CapitalizeFirstChar(routeName),
			PathProject: PathProject,
		}
		if arr, err := mkroute.New(); err != nil {
			log.Fatal(err)
		} else {
			for _, str := range arr {
				fmt.Printf("created %s\n", str)
			}

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
		if len(os.Args) < 3 {
			fmt.Println("Please Enter Model name")
			os.Exit(1)
		}
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
		db, err := database.MysqlConnect()
		if err != nil {
			log.Fatal(err)
		}
		model_name := os.Args[2]
		fmt.Printf("Do you want migrate %s Model? [y/N]:", model_name)
		Ans := "n"
		fmt.Scanf("%s\n", &Ans)
		if strings.ToLower(Ans) != "y" {
			log.Fatal("not migrate!")
		}

		err = database.Migrate(db, model_name)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s Migrate Success", model_name)

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
