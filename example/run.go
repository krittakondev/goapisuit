package main

import (
	"log"

	"github.com/krittakondev/goapisuit/v2"
	"github.com/krittakondev/goapisuit/v2/internal/routes"
	"github.com/krittakondev/goapisuit/v2/internal/setup"
)




func main(){
	suit, err := goapisuit.New("github.com/krittakondev/goapisuit/v2")
	if err != nil{
		log.Fatal(err)
	}

	suit.SetupRoutes(&routes.Route{})
	setup.GroupsSetup(suit)

	suit.Run()
}
