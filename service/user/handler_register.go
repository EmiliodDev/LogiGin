package user

import (
	"net/http"

	"github.com/EmiliodDev/LogiGin/service/auth"
	"github.com/EmiliodDev/LogiGin/types"
	"github.com/EmiliodDev/LogiGin/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleRegister(c *gin.Context) {
    var payload types.RegisterUserPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    if err := utils.Validate.Struct(payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    _, err := h.store.GetUserByEmail(payload.Email)
    if err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
        return
    }

    hashedPassword, err := auth.HashPassword(payload.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
        return
    }

    err = h.store.CreateUser(types.User{
        FirstName:  payload.FirstName,
        LastName:   payload.LastName,
        Email:      payload.Email,
        Password:   hashedPassword,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "registered user"})
}
