package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmiliodDev/todoAPI/service/auth"
	"github.com/EmiliodDev/todoAPI/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
    mock.Mock
}

func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
    args := m.Called(email)
    user, ok := args.Get(0).(*types.User)
    if !ok && args.Get(0) != nil {
        return nil, args.Error(1)
    }
    return user, args.Error(1)
}

func (m *MockUserStore) CreateUser(user types.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserStore) GetUserByID(id int) (*types.User, error) {
    args := m.Called(id)
    user, ok := args.Get(0).(*types.User)
    if !ok && args.Get(0) != nil {
        return nil, args.Error(1)
    }
    return user, args.Error(1)
}

func TestHandlerRegister(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    mockStore := new(MockUserStore)
    handler := NewHandler(mockStore)
    handler.RegisterRoutes(router.Group("/"))

    t.Run("successful registration", func(t *testing.T) {
        payload := types.RegisterUserPayload{
            FirstName:  "Emilio",
            LastName:   "Ortiz" ,
            Email:      "emilioortiz@example.com",
            Password:   "securepassword",
        }

        mockStore.On("GetUserByEmail", payload.Email).Return(types.User{}, errors.New("user not found"))
        mockStore.On("CreateUser", mock.AnythingOfType("types.User")).Return(nil).Run(func(args mock.Arguments) {
            user := args.Get(0).(types.User)
            assert.Equal(t, user.FirstName, payload.FirstName)
            assert.Equal(t, user.LastName, payload.LastName)
            assert.Equal(t, user.Email, payload.Email)
            assert.True(t, auth.ComparePasswords(user.Password, []byte(payload.Password)))
        })

        body, _ := json.Marshal(payload)
        req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()

        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
        mockStore.AssertExpectations(t)
    })

    t.Run("user already exists", func(t *testing.T) {
        payload := types.RegisterUserPayload{
            FirstName:  "Emilio",
            LastName:   "Ortiz",
            Email:      "emilioortiz@example.com",
            Password:   "securepassword",
        }

        existingUser := &types.User{
            FirstName:  "Emilio",
            LastName:   "Ortiz",
            Email:      "emilioortiz@example.com",
            Password:   "hashedpassword",
        }

        mockStore.On("GetUserByEmail", payload.Email).Return(existingUser, nil)

        body, _ := json.Marshal(payload)
        req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()

        router.ServeHTTP(resp, req)

        t.Logf("Response Code: %d, Body: %s", resp.Code, resp.Body.String())

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        mockStore.AssertExpectations(t)
    })


    t.Run("invalid payload", func(t *testing.T) {
        body := []byte(`{"firstName": "", "lastName": "", "email": "invalid-email", "password": ""}`)
        req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()

        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
    })
}
