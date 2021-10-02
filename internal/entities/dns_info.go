package entities

import "gorm.io/gorm"

type DNSInfo struct {
	gorm.Model
	Status      string
	CnameRecord string
	ARecords    string
	NSRecords   string
}
