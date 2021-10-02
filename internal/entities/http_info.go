package entities

import "gorm.io/gorm"

type HTTPInfo struct {
	gorm.Model
	HTTPStatusCode  int
	HTTPSStatusCode int
	Title           string
	WebServer       string
	ContentType     string
}
