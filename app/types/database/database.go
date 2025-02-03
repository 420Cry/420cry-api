package types

import "gorm.io/gorm"

type Database struct {
	DB *gorm.DB
}
