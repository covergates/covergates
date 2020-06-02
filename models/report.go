package models

import "github.com/jinzhu/gorm"

type Report struct {
	gorm.Model
	Data   []byte
	Type   string
	Branch string `gorm:"index"`
	Tag    string `gorm:"index"`
	Commit string `gorm:"index"`
}
