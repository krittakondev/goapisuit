package main

import (
	"log"

	"github.com/krittakondev/goapisuit"
	routes "github.com/krittakondev/goapisuit/internal/api"
)




func main(){
	suit, err := goapisuit.New("github.com/krittakondev/goapisuit")
	if err != nil{
		log.Fatal(err)
	}
	suit.Run(&routes.Route{Suit: suit})
}
