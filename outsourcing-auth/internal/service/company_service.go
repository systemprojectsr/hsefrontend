package service

import (
	"core/internal"
	"core/internal/api"
	"core/internal/database"
	"core/internal/database/repository"
	"core/internal/security"
	"errors"
	"strconv"
)

type CompanyService interface {
	Signup(request *api.UserCompanyRegister) (database.CompanyDB, error)
	GetCompany(id uint) (database.CompanyDB, error)
	Login(request *api.GeneralAuth) (dbUser database.CompanyDB, jwtToken string, err error)
	AccessByToken(request *api.TokenAccess) (*api.ResponseSuccessAccess, database.CompanyDB, error)
	CreateCard(request *api.TokenCreateCard) (error, database.Card)
	ListCard(request *api.TokenListCard, limit string, page string) (error, []database.Card)
	DeleteCard(request *api.TokenDeleteCard) (error, bool)
}

type companyService struct {
	repository repository.CompanyRepository
}

func (service *companyService) Signup(request *api.UserCompanyRegister) (database.CompanyDB, error) {
	exists, existsCompany, _ := service.repository.ExistsByEmail(request.Email)
	if exists || existsCompany {
		return database.CompanyDB{}, errors.New("email already exists")
	}

	company := &database.CompanyDB{
		CompanyName:   request.CompanyName,
		FullName:      request.FullName,
		PositionAgent: request.PositionAgent,
		IDCompany:     request.IDCompany,
		Email:         request.Email,
		Phone:         request.Phone,
		Address:       request.Address,
		TypeService:   request.TypeService,
		PasswordHash:  request.PasswordHash,
		Photo:         request.Photo,
		Documents:     request.Documents,
		Stars:         5,
		Type:          "company",
	}
	service.repository.Save(company)
	return *company, nil
}

func (service *companyService) GetCompany(id uint) (database.CompanyDB, error) {
	company, err := service.repository.GetByID(id)
	if err != nil {
		return database.CompanyDB{}, err
	}
	return *company, nil
}

func (service *companyService) Login(request *api.GeneralAuth) (dbUser database.CompanyDB, jwtToken string, err error) {
	dbUser, err = service.repository.CheckPassword(request.GeneralLogin.GeneralLoginAttributes.Email, request.GeneralLogin.GeneralLoginAttributes.PasswordHash)
	if err != nil {
		return dbUser, "", err
	}
	jwtToken = security.CreateToken(dbUser.Type == "company", dbUser.ID, internal.LifeTimeJWT)
	if jwtToken == "" {
		return dbUser, "", errors.New("the created jwt was faulty")
	}
	return dbUser, jwtToken, nil
}

func (service *companyService) AccessByToken(request *api.TokenAccess) (*api.ResponseSuccessAccess, database.CompanyDB, error) {
	result, tokenStructure := security.CheckToken(request.User.Login.Token)
	company, err := service.GetCompany(uint(tokenStructure["accessID"].(float64)))
	if err != nil {
		return nil, company, err
	}

	if result {
		response := api.ResponseSuccessAccess{
			StatusResponse: &internal.StatusResponse{Status: "success"},
			ResponseUser: &api.ResponseUser{
				ID:    company.ID,
				Token: request.User.Login.Token,
				Type:  company.Type,
			},
		}
		return &response, company, nil
	} else {
		response := api.ResponseSuccessAccess{
			StatusResponse: &internal.StatusResponse{Status: "success"},
			ResponseUser: &api.ResponseUser{
				ID:    company.ID,
				Token: "expired",
				Type:  company.Type,
			},
		}
		return &response, company, nil
	}
}

func (service *companyService) CreateCard(request *api.TokenCreateCard) (error, database.Card) {
	result, tokenStructure := security.CheckToken(request.User.Login.Token)
	company, err := service.GetCompany(uint(tokenStructure["accessID"].(float64)))
	if err != nil {
		return err, database.Card{}
	}
	if result {
		card := database.Card{
			Title:       request.Card.Title,
			Description: request.Card.Description,
			CompanyID:   company.ID,
		}
		resp := service.repository.SaveCard(&card)
		if resp {
			return nil, card
		} else {
			return errors.New("resp in CreateCard()"), database.Card{}
		}
	} else {
		return err, database.Card{}
	}
}

func (service *companyService) DeleteCard(request *api.TokenDeleteCard) (error, bool) {
	result, tokenStructure := security.CheckToken(request.User.Login.Token)
	_, err := service.GetCompany(uint(tokenStructure["accessID"].(float64)))
	if err != nil {
		return err, false
	}
	if result {
		resp := service.repository.DeleteCard(request.Card.ID)
		if resp {
			return nil, resp
		} else {
			return err, false
		}
	} else {
		return err, false
	}
}

func (service *companyService) ListCard(request *api.TokenListCard, limit string, page string) (error, []database.Card) {
	result, tokenStructure := security.CheckToken(request.TokenAccess.User.Login.Token)
	company, err := service.GetCompany(uint(tokenStructure["accessID"].(float64)))
	if err != nil {
		return err, []database.Card{}
	}
	if result {
		limitI, err := strconv.Atoi(limit)
		if err != nil && limit != "" {
			return err, []database.Card{}
		}
		pageI, err := strconv.Atoi(page)
		if err != nil && page != "" {
			return err, []database.Card{}
		}
		service.repository.PreloadDB("Cards", &company, limitI, pageI)
		return nil, company.Cards
	} else {
		return nil, []database.Card{}
	}
}

func NewCompanyService(repository repository.CompanyRepository) CompanyService {
	return &companyService{
		repository: repository,
	}
}
