package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ClientDB struct {
	gorm.Model
	ID           uint `gorm:"primaryKey;autoIncrement"`
	FullName     string
	Email        string
	Phone        string
	PasswordHash string
	Photo        string
	Type         string
	Permissions  pq.StringArray `gorm:"type:text[]"`
}

type CompanyDB struct {
	gorm.Model
	ID            uint `gorm:"primaryKey;autoIncrement"`
	CompanyName   string
	FullName      string
	PositionAgent string
	IDCompany     string
	Email         string
	Phone         string
	Address       string
	TypeService   string
	PasswordHash  string
	Photo         string
	Documents     pq.StringArray `gorm:"type:text[]"`
	Stars         float64
	Type          string
	Permissions   pq.StringArray `gorm:"type:text[]"`
	Cards         []Card         `gorm:"foreignKey:CompanyID"`
}

type Card struct {
	gorm.Model
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Title       string
	Description string
	CompanyID   uint
	Company     CompanyDB `gorm:"foreignKey:CompanyID"`
}
