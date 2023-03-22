package entities

import "gorm.io/gorm"

type Subdomain struct {
	gorm.Model
	Domain     string
	ProgramID  *int `json:"-"`
	Program    *Program
	Tags       []*Tag  `gorm:"many2many:subdomain_tags;"`
	Techs      []*Tech `gorm:"many2many:subdomain_techs;"`
	DNSInfoID  *int    `json:"-"`
	DNSInfo    *DNSInfo
	HTTPInfoID *int `json:"-"`
	HTTPInfo   *HTTPInfo
}
