package maketemplate

import (
	"log"
	"os"
	"text/template"
)


type MakeRoute struct{
	Name string
}

func (mr  *MakeRoute) New() error{
	tmpl, err := template.ParseFiles("pkg/maketemplate/tmpl/route.go.tmpl")
	if err != nil{
		return err
	}
	createPath := "internal/api/"+mr.Name+".go"

	if _, err := os.Stat(createPath); err == nil {
		log.Fatal(createPath+ " is Exist")
	}


	file, err := os.OpenFile(createPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	

	if err := tmpl.Execute(file, mr); err != nil{
		return err
	}

	return nil
}
