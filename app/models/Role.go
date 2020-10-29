package models

import "gorm.io/gorm"

// Role model
type Role struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name" xml:"name" form:"name" query:"name"`
	Description string `gorm:"type:varchar(100);" json:"description"`
}
