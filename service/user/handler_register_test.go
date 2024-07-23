package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmiliodDev/LogiGin/service/auth"
	"github.com/EmiliodDev/LogiGin/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerRegister(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockStore := new(MockUserStore)
    handler := NewHandler(mockStore)

    t.Run("successful registration", func(t *testing.T) {
        payload := types.RegisterUserPayload{
            FirstName:  "Emilio",
            LastName:   "Ortiz" ,
            Email:      "emilioortiz@example.com",
            Password:   "securepassword",
        }

        mockStore.On("GetUserByEmail", payload.Email).Return(nil, errors.New("user not found"))
        mockStore.On("CreateUser", mock.AnythingOfType("types.User")).Return(nil).Run(func(args mock.Arguments) {
            user := args.Get(0).(types.User)
            assert.Equal(t, user.FirstName, payload.FirstName)
            assert.Equal(t, user.LastName, payload.LastName)
            assert.Equal(t, user.Email, payload.Email)
            assert.True(t, auth.ComparePasswords(user.Password, []byte(payload.Password)))
        })

        body, err := json.Marshal(payload)
        if err != nil {
            t.Fatal(err)
        }
        req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
        if err != nil {
            t.Fatal(err)
        }
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()
        
        router := gin.Default()

        router.POST("/register", handler.handleRegister)
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
        mockStore.AssertExpectations(t)
    })

    t.Run("user already exists", func(t *testing.T) {
        payload := types.RegisterUserPayload{
            FirstName:  "Emilio",
            LastName:   "Ortiz" ,
            Email:      "emilioortiz@exampl.com",
            Password:   "securepassword",
        }

        existingUser := &types.User{
            ID:             1,
            FirstName:      "Emilio",
            LastName:       "Ortiz",
            Email:          "emilioortiz@example.com",
            Password:       "hashedpassword",
            CreatedAt:      time.Time{},
        }

        mockStore.On("GetUserByEmail", payload.Email).Return(existingUser, nil)

        body, err := json.Marshal(payload)
        if err != nil {
            t.Fatal(err)
        }

        req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
        if err != nil {
            t.Fatal(err)
        }

        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()
        
        router := gin.Default()

        router.POST("/register", handler.handleRegister)
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        mockStore.AssertExpectations(t)
    })

    t.Run("invalid payload", func(t *testing.T) {
        body := []byte(`{"firstName": "", "lastName": "", "email": "invalid-email", "password": ""}`)
        req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        resp := httptest.NewRecorder()

        router := gin.Default()
        router.POST("/register", handler.handleRegister)
        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
    })
}

type MockUserStore struct {
    mock.Mock
}

func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
    args := m.Called(email)
    
    if user, ok := args.Get(0).(*types.User); ok {
        return user, args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *MockUserStore) CreateUser(user types.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserStore) GetUserByID(id int) (*types.User, error) {
    args := m.Called(id)
    user, ok := args.Get(0).(*types.User)
    if !ok {
        return nil, args.Error(1)
    }
    return user, nil
}

