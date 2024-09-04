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
	tmplRoute, err := template.ParseFiles("pkg/maketemplate/tmpl/route.go.tmpl")
	if err != nil{
		return err
	}
	tmplModel, err := template.ParseFiles("pkg/maketemplate/tmpl/model.go.tmpl")
	if err != nil{
		return err
	}
	createPathRoute := "internal/api/"+mr.Name+".go"
	createPathModel := "internal/models/"+mr.Name+".go"

	if _, err := os.Stat(createPathRoute); err == nil{
		log.Fatal(createPathRoute+ " is Exist")
	}
	if _, err := os.Stat(createPathModel); err == nil{
		log.Fatal(createPathModel+ " is Exist")
	}

	file, err := os.OpenFile(createPathRoute, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmplRoute.Execute(file, mr); err != nil{
		return err
	}

	file, err = os.OpenFile(createPathModel, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmplModel.Execute(file, mr); err != nil{
		return err
	}

	return nil
}
