package main

import (
	"log"

	"github.com/krittakondev/goapisuit"
	"github.com/krittakondev/goapisuit/internal/routes"
)




func main(){
	suit, err := goapisuit.New("github.com/krittakondev/goapisuit")
	if err != nil{
		log.Fatal(err)
	}
	suit.Run(&routes.Route{Suit: suit})
}
