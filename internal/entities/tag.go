package entities

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name       string
	Subdomains []*Subdomain `gorm:"many2many:subdomain_tags;"`
}
