package user

import (
	"fmt"
	"net/http"

	"github.com/EmiliodDev/LogiGin/config"
	"github.com/EmiliodDev/LogiGin/service/auth"
	"github.com/EmiliodDev/LogiGin/types"
	"github.com/EmiliodDev/LogiGin/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) handleLogin(c *gin.Context) {
    var payload types.LoginUserPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    if err := utils.Validate.Struct(payload); err != nil {
        errors := err.(validator.ValidationErrors)
        c.JSON(http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
        return
    }

    u, err := h.store.GetUserByEmail(payload.Email)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "not found, invalid email or password"})
        return
    }
    
    if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
        return
    }

    secret := []byte(config.Envs.JWTSecret)
    token, err := auth.CreateJWT(secret, u.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, map[string]string{"token": token})
}


