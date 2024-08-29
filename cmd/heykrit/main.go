package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/krittakondev/goapisuit/pkg/maketemplate"
)

func capitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

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
		routeName := capitalizeFirstChar(os.Args[2])

		mkroute := &maketemplate.MakeRoute{
			Name: routeName,
		}
		if err := mkroute.New(); err != nil{
			log.Fatal(err)
		}

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
