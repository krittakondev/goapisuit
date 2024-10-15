package models

import "gorm.io/gorm"


  // Edit this models
type Users struct {
	gorm.Model
	Name        string `json:"name,omitempty" gorm:"index"`
	Username    string `json:"username,omitempty" gorm:"unique"`
	PasswordEnc string `json:"peassword_enc,omitempty"`
}
