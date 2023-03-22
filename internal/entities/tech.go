package entities

import "gorm.io/gorm"

type Tech struct {
	gorm.Model
	Name       string
	Subdomains []*Subdomain `gorm:"many2many:subdomain_techs;"`
}

func (Tech) TableName() string {
	return "techs"
}
