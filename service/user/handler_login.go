package user

import (
	"fmt"
	"net/http"

	"github.com/EmiliodDev/todoAPI/service/auth"
	"github.com/EmiliodDev/todoAPI/types"
	"github.com/EmiliodDev/todoAPI/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) handleLogin(c *gin.Context) {
    var user types.LoginUserPayload
    if err := utils.ParseJSON(c, &user); err != nil {
        utils.WriteError(c, http.StatusBadRequest, err)
        return
    }

    if err := utils.Validate.Struct(user); err != nil {
        errors := err.(validator.ValidationErrors)
        utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
        return
    }

    u, err := h.store.GetUserByEmail(user.Email)
    if err != nil {
        utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("not found invalid email or password"))
        return
    }
    
    if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
        utils.WriteError(c, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
        return
    }

    utils.WriteJSON(c, http.StatusOK, map[string]string{"token": ""})
}


