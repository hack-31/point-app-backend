package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct{}

func NewHealthCheckHandler() *HealthCheck {
	return &HealthCheck{}
}

func (hc *HealthCheck) ServeHTTP(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
