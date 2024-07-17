package user

import (
	"github.com/EmiliodDev/todoAPI/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
    store   types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
    return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/register", h.handleRegister)
    router.POST("/login", h.handleLogin)
}

func (h *Handler) handleLogin(c *gin.Context) {
    
}

func (h *Handler) handleRegister(c *gin.Context) {

}
