package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handleGetUser(c *gin.Context) {
    str := c.Param("userID")

    userID, err := strconv.Atoi(str)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid user ID"})
        return
    }

    user, err := h.store.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, user)
}
