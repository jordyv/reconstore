package entities

import "gorm.io/gorm"

type Program struct {
	gorm.Model
	Name        string
	Slug        string
	Platform    string
	Private     bool
	HasBounties bool
}
