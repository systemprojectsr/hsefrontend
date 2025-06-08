package main

import (
	"core/internal"
	"core/internal/api"
	"core/internal/controller"
	"core/internal/database"
	"core/internal/database/repository"
	"core/internal/security"
	"core/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type Entity interface {
	database.ClientDB | database.CompanyDB
}

func main() {
	r := gin.Default()

	err := internal.InitEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitialiseDB(&database.DbConfig{
		User:     internal.PostgresUser,
		Password: internal.PostgresPassword,
		DbName:   internal.PostgresDB,
		Host:     internal.PostgresHost,
		Port:     internal.PostgresPort,
		Schema:   "account",
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&database.ClientDB{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.CompanyDB{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&database.Card{})
	if err != nil {
		panic(err)
	}

	clientRepository := repository.NewClientRepository(db)
	clientService := service.NewClientService(clientRepository)
	clientController := controller.NewClientController(clientService)
	companyRepository := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepository)
	companyController := controller.NewCompanyController(companyService)

	v1 := r.Group("v1")
	{
		v1.POST("/login", func(c *gin.Context) {
			request := &api.GeneralAuth{}
			if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
				api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
				return
			}
			resultClient, _, _ := clientRepository.ExistsByEmail(request.GeneralLogin.GeneralLoginAttributes.Email)
			if resultClient {
				clientController.Login(c, request)
			} else {
				companyController.Login(c, request)
			}
		})
		accountGroup := v1.Group("account")
		{
			accountGroup.POST("/", func(c *gin.Context) {
				request := &api.TokenAccess{}
				if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
					api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
					return
				}
				ok, mapClaims := security.CheckToken(request.User.Login.Token)
				if mapClaims == nil {
					api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
					return
				}
				if ok {
					isCompany := mapClaims["isCompany"].(bool)
					if isCompany {
						companyController.GetAccount(c, request)
					} else {
						clientController.GetAccount(c, request)
					}
				} else {
					api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
					return
				}
			})
			cardGroup := accountGroup.Group("card")
			{
				cardGroup.POST("/create", func(c *gin.Context) {
					request := &api.TokenCreateCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							companyController.CreateCard(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
				cardGroup.POST("/list", func(c *gin.Context) {
					request := &api.TokenListCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					limit := c.Query("limit")
					page := c.Query("page")
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							companyController.ListCard(c, request, limit, page)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
				cardGroup.POST("/delete", func(c *gin.Context) {
					request := &api.TokenDeleteCard{}
					if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
						api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
						return
					}
					ok, mapClaims := security.CheckToken(request.User.Login.Token)
					if mapClaims == nil {
						api.GetErrorJSON(c, http.StatusBadRequest, "The token is invalid")
						return
					}
					if ok {
						isCompany := mapClaims["isCompany"].(bool)
						if isCompany {
							companyController.DeleteCard(c, request)
						} else {
							api.GetErrorJSON(c, http.StatusForbidden, "You're not a company")
							return
						}
					} else {
						api.GetErrorJSON(c, http.StatusForbidden, "The token had expired")
						return
					}
				})
			}
		}
		registerGroup := v1.Group("register")
		{
			registerGroup.POST("/client", func(c *gin.Context) {
				clientController.Signup(c)
			})
			registerGroup.POST("/company", func(c *gin.Context) {
				companyController.Signup(c)
			})
		}
	}

	r.Run()
}
