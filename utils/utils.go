package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(c *gin.Context, payload any) error {
    if c.Request.Body == nil {
        return fmt.Errorf("missing request body")
    }
    
    return c.ShouldBindJSON(payload)
}

func WriteJSON(c *gin.Context, status int, v any) {
   c.JSON(status, v)
}

func WriteError(c *gin.Context, status int, err error) {
    WriteJSON(c, status, map[string]string{"error": err.Error()})
}
