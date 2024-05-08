package link

import "gorm.io/gorm"

type Link struct {
	gorm.Model
	Path string
	URL  string
}
