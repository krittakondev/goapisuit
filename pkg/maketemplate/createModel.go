package maketemplate

import (
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
	if err := tmpl.Execute(os.Stdout, mr); err != nil{
		return err
	}

	return nil
}
