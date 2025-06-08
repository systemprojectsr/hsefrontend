package controller

import (
	"core/internal"
	"core/internal/api"
	"core/internal/security"
	"core/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ClientController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context, request *api.GeneralAuth)
	GetAccount(c *gin.Context, request *api.TokenAccess)
}

type clientController struct {
	service service.ClientService
}

func (controller clientController) Signup(c *gin.Context) {
	request := &api.ClientRegister{}
	if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
		return
	}
	client, err := controller.service.Signup(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusPreconditionFailed, err.Error())
		return
	}
	tokenGenerated := security.CreateToken(false, client.ID, 60)
	c.JSON(http.StatusOK, api.ResponseSuccessAccess{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		ResponseUser: &api.ResponseUser{
			ID:    client.ID,
			Token: tokenGenerated,
			Type:  client.Type,
		},
	})
}

func (controller clientController) Login(c *gin.Context, request *api.GeneralAuth) {
	dbUser, jwtToken, err := controller.service.Login(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "the created jwt was faulty")
		return
	}
	c.JSON(http.StatusOK, api.ResponseSuccessAccess{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		ResponseUser: &api.ResponseUser{
			ID:    dbUser.ID,
			Token: jwtToken,
			Type:  dbUser.Type,
		},
	})
}

func (controller clientController) GetAccount(c *gin.Context, request *api.TokenAccess) {
	response, user, err := controller.service.AccessByToken(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "The token is incorrect")
		return
	}
	if response.ResponseUser.Token == "expired" {
		api.GetErrorJSON(c, http.StatusForbidden, "The token has expired")
		return
	}
	c.JSON(http.StatusOK, api.ResponseAccount{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		User: struct {
			Account api.AccountInfo `json:"account"`
		}{Account: api.AccountInfo{
			ID:       user.ID,
			FullName: user.FullName,
			Phone:    user.Phone,
			Photo:    user.Photo,
			Token:    response.ResponseUser.Token,
			Type:     user.Type,
		}},
	})
}

func NewClientController(service service.ClientService) ClientController {
	return &clientController{
		service: service,
	}
}
