package api

import (
	"core/internal"
	"github.com/gin-gonic/gin"
)

func GetErrorJSON(c *gin.Context, status int, description string) {
	response := internal.InfoResponse{
		StatusResponse: &internal.StatusResponse{Status: "error"},
		Description:    description,
	}
	c.JSON(status, response)
}
