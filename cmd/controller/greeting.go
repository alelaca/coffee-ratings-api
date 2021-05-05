package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Greeting(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello world!")
}
