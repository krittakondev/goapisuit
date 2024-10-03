package main

import (
	"log"

	"github.com/krittakondev/goapisuit"
	"github.com/krittakondev/goapisuit/internal/routes"
	"github.com/krittakondev/goapisuit/internal/setup"
)




func main(){
	suit, err := goapisuit.New("github.com/krittakondev/goapisuit")
	if err != nil{
		log.Fatal(err)
	}

	suit.SetupRoutes(&routes.Route{})
	setup.GroupsSetup(suit)

	suit.Run()
}
