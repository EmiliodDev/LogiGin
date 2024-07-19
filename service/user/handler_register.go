package user

import (
	"fmt"
	"net/http"

	"github.com/EmiliodDev/todoAPI/service/auth"
	"github.com/EmiliodDev/todoAPI/types"
	"github.com/EmiliodDev/todoAPI/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleRegister(c *gin.Context) {
    var user types.RegisterUserPayload
    if err := utils.ParseJSON(c, &user); err != nil {
        utils.WriteError(c, http.StatusBadRequest, err)
        return
    }

    if err := utils.Validate.Struct(user); err != nil {
        utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", err))
        return
    }

    _, err := h.store.GetUserByEmail(user.Email)
    if err == nil {
        utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", user.Email))
        return
    }

    hashedPassword, err := auth.HashPassword(user.Password)
    if err != nil {
        utils.WriteError(c, http.StatusInternalServerError, err)
        return
    }

    err = h.store.CreateUser(types.User{
        FirstName: user.FirstName,
        LastName: user.LastName,
        Email: user.Email,
        Password: hashedPassword,
    })
    if err != nil {
        utils.WriteError(c, http.StatusInternalServerError, err)
        return
    }

    utils.WriteJSON(c, http.StatusCreated, nil)
}
