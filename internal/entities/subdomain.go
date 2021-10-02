package entities

import "gorm.io/gorm"

type Subdomain struct {
	gorm.Model
	Domain    string
	ProgramID int
	Program   Program
	Tags      []*Tag `gorm:"many2many:subdomain_tags;"`
	DNSInfoID int
	DNSInfo   DNSInfo
}
