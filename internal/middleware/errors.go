package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
)

func JSONAppErrorReporter() gin.HandlerFunc {
	return HadleError(gin.ErrorTypeAny)
}

func HadleError(errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Skip if no errors
		if c.Errors.Last() == nil {
			return
		}
		// public errors
		err := c.Errors.Last()
		if err == nil {
			return
		}

		c.JSON(-1, domain.ErrorResponse{
			Code:    c.Writer.Status(),
			Message: err.Error(),
		})
		c.Abort()
	}
}
