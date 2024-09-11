package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/krittakondev/goapisuit/internal/database"
	"github.com/krittakondev/goapisuit/pkg/maketemplate"
	"github.com/krittakondev/goapisuit/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass a command")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		// template := &maketemplate.Template{}
		fmt.Println("create init")
	
	case "make":
		if len(os.Args) < 3 {
			fmt.Println("Please Enter Route name")
			os.Exit(1)
		}
		routeName := os.Args[2]

		mkroute := &maketemplate.MakeRoute{
			Name: utils.CapitalizeFirstChar(routeName),
		}
		if arr, err := mkroute.New(); err != nil {
			log.Fatal(err)
		}else{
			for _, str := range arr{
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
		fmt.Scanf("%s\n",&Ans)
		if strings.ToLower(Ans) !=  "y"{
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
