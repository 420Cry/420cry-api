// Package middleware provides centralized error handling for the application.
package middleware

import (
	"log"
	"net/http"

	app_errors "cry-api/app/types/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler provides centralized error handling middleware
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err.Err)
		}
	}
}

// handleError processes different types of errors and returns appropriate responses
func handleError(c *gin.Context, err error) {
	log.Printf("Error occurred: %v", err)

	switch e := err.(type) {
	case *app_errors.ValidationError:
		c.JSON(e.Code, gin.H{
			"error": e.Message,
			"field": e.Field,
			"value": e.Value,
		})
	case *app_errors.NotFoundError:
		c.JSON(e.Code, gin.H{
			"error": e.Message,
		})
	case *app_errors.UnauthorizedError:
		c.JSON(e.Code, gin.H{
			"error": e.Message,
		})
	case *app_errors.ConflictError:
		c.JSON(e.Code, gin.H{
			"error": e.Message,
		})
	case *app_errors.InternalServerError:
		c.JSON(e.Code, gin.H{
			"error": e.Message,
		})
	case *app_errors.AppError:
		c.JSON(e.Code, gin.H{
			"error": e.Message,
		})
	default:
		// Generic error handling
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
}

// AbortWithError aborts the request with a specific error
func AbortWithError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
