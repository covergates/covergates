package models

import "github.com/jinzhu/gorm"

type Repo struct {
	gorm.Model
	URL      string `gorm:"unique_index;not null"`
	ReportID string
	Name     string `gorm:"unique_index:repo_name"`
	SCM      string `gorm:"unique_index:repo_name"`
}
