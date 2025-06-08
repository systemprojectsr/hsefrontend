package repository

import (
	"core/internal/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	Save(company *database.CompanyDB)
	ExistsByEmail(email string) (exists bool, existsClient bool, companyDB database.CompanyDB)
	GetByID(id uint) (*database.CompanyDB, error)
	CheckPassword(email string, passwordHash string) (database.CompanyDB, error)
	SaveCard(card *database.Card) bool
	PreloadDB(name string, company *database.CompanyDB, limit int, page int)
	DeleteCard(id int) bool
}

type companyRepository struct {
	db *gorm.DB
}

func (repository *companyRepository) PreloadDB(name string, company *database.CompanyDB, limit int, page int) {
	query := repository.db

	if limit > 0 {
		offset := page * limit
		query = query.Preload(name, func(db *gorm.DB) *gorm.DB {
			return db.Limit(limit).Offset(offset)
		})
	} else {
		query = query.Preload(name)
	}

	query.First(&company, company.ID)
}

func (repository *companyRepository) Save(company *database.CompanyDB) {
	repository.db.Save(company)
}

func (repository *companyRepository) SaveCard(card *database.Card) bool {
	resp := repository.db.Save(card)
	if resp != nil {
		return true
	} else {
		return false
	}
}

func (repository *companyRepository) DeleteCard(id int) bool {
	resp := repository.db.Delete(&database.Card{}, id)
	if resp != nil {
		return true
	} else {
		return false
	}
}

func (repository *companyRepository) ExistsByEmail(email string) (exists bool, existsClient bool, companyDB database.CompanyDB) {
	var company database.CompanyDB
	var client database.ClientDB
	result := repository.db.Model(&database.CompanyDB{}).Where("email = ?", email).First(&company)
	exists = !errors.Is(result.Error, gorm.ErrRecordNotFound)
	if exists == false {
		result = repository.db.Model(&database.ClientDB{}).Where("email = ?", email).First(&client)
		existsClient = !errors.Is(result.Error, gorm.ErrRecordNotFound)
	}
	return exists, existsClient, company
}

func (repository *companyRepository) GetByID(id uint) (*database.CompanyDB, error) {
	var company database.CompanyDB

	result := repository.db.Model(&database.CompanyDB{}).Where("id = ?", id).First(&company)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &company, nil
}

func (repository *companyRepository) CheckPassword(email string, passwordHash string) (database.CompanyDB, error) {
	ok, _, dbUser := repository.ExistsByEmail(email)
	if ok {
		if dbUser.PasswordHash == passwordHash {
			return dbUser, nil
		}
		return database.CompanyDB{}, errors.New("bad password hash")
	}
	return database.CompanyDB{}, errors.New("user not found")
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	var repository CompanyRepository

	repository = &companyRepository{
		db: db,
	}

	return repository
}
