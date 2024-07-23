package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func GetTokenFromRequest(c *gin.Context) string {
    tokenAuth := c.GetHeader("Authorization")
    tokenQuery := c.Query("token")

    if tokenAuth != "" {
        return tokenAuth
    }

    if tokenQuery != "" {
        return tokenQuery
    }

    return ""
}
