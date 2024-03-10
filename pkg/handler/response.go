package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"testApi"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type getAllGroupsResponse struct {
	Data []testApi.Group `json:"data"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
