package featA

import "github.com/gin-gonic/gin"

type AvailabilityHandler interface {
	CheckAvailability(c *gin.Context)
}