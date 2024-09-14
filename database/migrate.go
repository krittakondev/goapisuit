package database

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, model_name string) error{
	
	switch strings.ToLower(model_name){

	default:
		return errors.New("not found model")
	}
}
