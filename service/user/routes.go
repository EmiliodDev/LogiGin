package user

import (
	"github.com/EmiliodDev/LogiGin/service/auth"
	"github.com/EmiliodDev/LogiGin/types"
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

    router.GET("/users/{userID}", auth.WithJWTAuth(h.handleGetUser, h.store))
}
