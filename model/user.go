package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `mapstructure:"username"`
	PassWord string `mapstructure:"pwd"`
}

func (User) TableName() string {
	return "user"
}
