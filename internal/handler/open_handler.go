package handler

import (
	"net/http"

	"github.com/ZorinIvanA/http-gate-control/internal/service"
	"github.com/gin-gonic/gin"
)

func OpenHandler(gateSvc service.GateService) gin.HandlerFunc {
	return func(c *gin.Context) {
		opened, err := gateSvc.ProcessOpen()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if opened {
			c.JSON(http.StatusOK, gin.H{"status": "opened"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "blocked"})
		}
	}
}
