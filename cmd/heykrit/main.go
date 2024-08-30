package main

import (
	"fmt"
	"log"
	"os"

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
	case "make":
		if len(os.Args) < 3 {
			fmt.Println("Please Enter Route name")
			os.Exit(1)
		}
		routeName := os.Args[2]

		mkroute := &maketemplate.MakeRoute{
			Name:  utils.CapitalizeFirstChar(routeName),
		}
		if err := mkroute.New(); err != nil{
			log.Fatal(err)
		}

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
