package api

import (
	"core/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMockStandard(c *gin.Context) {
	response := internal.InfoResponse{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		Description:    "This is a mock",
	}
	c.JSON(http.StatusOK, response)
}
