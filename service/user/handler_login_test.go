package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EmiliodDev/todoAPI/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandlerLogin(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mocksStore := new(MockUserStore)
    handler := NewHandler(mocksStore)

    t.Run("successful login", func(t *testing.T) {
        payload := types.LoginUserPayload{
            Email:      "test@test.com",
            Password:   "securepassword",
        }

        existingUser := &types.User{
            ID:         1,
            FirstName:  "Emilio",
            LastName:   "Ortiz",
            Email:      "test@test.com",
            Password:   "hashedpassword",
            CreatedAt:  time.Time{},
        }

        mocksStore.On("GetUserByEmail", payload.Email).Return(existingUser, nil)

        body, err := json.Marshal(payload)
        if err != nil {
            t.Fatal(err)
        }

        req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
        if err != nil {
            t.Fatal(err)
        }

        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()

        router := gin.Default()

        router.POST("/login", handler.handleLogin)
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
        mocksStore.AssertExpectations(t)
    })
}
