package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Serve() error {
	router := gin.Default()
	router.GET("/health", getHealth)
	
	if err := router.Run("0.0.0.0:8080"); err != nil {
		return err
	}
	return nil
}

func getHealth(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
