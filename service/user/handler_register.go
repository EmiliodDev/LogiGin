package user

import (
	"fmt"
	"net/http"

	"github.com/EmiliodDev/todoAPI/types"
	"github.com/EmiliodDev/todoAPI/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleRegister(c *gin.Context) {
    var user types.RegisterUserPayload
    if err := utils.ParseJSON(c.Request, &user); err != nil {
        utils.WriteError(c.Writer, http.StatusBadRequest, err)
        return
    }

    if err := utils.Validate.Struct(user); err != nil {
        utils.WriteError(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload %v", err))
        return
    }



}
