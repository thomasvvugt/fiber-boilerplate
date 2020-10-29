package models

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Password string `json:"-" xml:"-" form:"-" query:"-"`
	Email string
	RoleID uint `gorm:"column:role_id" json:"role_id"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
