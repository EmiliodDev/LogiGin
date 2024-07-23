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
)

func TestHandlerLogin(t *testing.T) {
    gin.SetMode(gin.TestMode)

    tests := []struct {
        name                string
        payload             interface{}
        setupMocks          func(*MockUserStore)
        expectedStatus      int
        expectedResponse    string
    }{
        {
            name:               "invalid JSON payload",
            payload:            `{
                                    "email": "invalid-email",
                                    "password": "securepassword"
                                }`,
            setupMocks:         func(mus *MockUserStore) {},
            expectedStatus:     http.StatusBadRequest,
            expectedResponse:   `{"error":"invalid payload"}`,
        },
        {
            name:               "validation error",
            payload:            types.LoginUserPayload{
                                    Email: "invalid-email",
                                    Password: "securepassword",
                                },
            setupMocks:         func(mus *MockUserStore) {},
            expectedStatus:     http.StatusBadRequest,
            expectedResponse:   `{}`,
        },
        {
            name:               "user not found",
            payload:            types.LoginUserPayload{
                                    Email: "user@example.com",
                                    Password: "securepassword",
                                },
            setupMocks:         func(mus *MockUserStore) {
                                    mus.On("GetUserByEmail", "user@example.com").Return(nil, errors.New("user not found"))
                                },
            expectedStatus:     http.StatusBadRequest,
            expectedResponse:   `{"error":"not found, invalid email or password"}`,
        },
        {
            name:               "incorrect password",
            payload:            types.LoginUserPayload{
                                    Email:      "user@example.com",
                                    Password:   "wrongpassword",
                                },
            setupMocks:         func(mus *MockUserStore) {
                                    hashedPassword, _ := auth.HashPassword("correctpassword") 
                                    mus.On("GetUserByEmail", "user@example.com").Return(&types.User{Email:"user@example.com", Password: hashedPassword}, nil)
                                },
            expectedStatus:     http.StatusBadRequest,
            expectedResponse:   `{"error":"invalid email or password"}`,
        },
        {
            name:               "successful login",
            payload:            types.LoginUserPayload{
                                    Email:      "user@example.com",
                                    Password:   "correctpassword",
                                },
            setupMocks:         func(mus *MockUserStore) {
                                    hashedPassword, _ := auth.HashPassword("correctpassword") 
                                    mus.On("GetUserByEmail", "user@example.com").Return(&types.User{Email: "user@example.com", Password: hashedPassword}, nil)
                                },
            expectedStatus:     http.StatusOK,
            expectedResponse:   `{"token":""}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T)  {
            mockStore := new(MockUserStore)

            tt.setupMocks(mockStore)

            h := NewHandler(mockStore)


            recorder := httptest.NewRecorder()
            context, _ := gin.CreateTestContext(recorder)

            payloadBytes, _ := json.Marshal(tt.payload)
            req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payloadBytes))
            context.Request = req

            h.handleLogin(context)

            assert.Equal(t, tt.expectedStatus, recorder.Code)
            assert.JSONEq(t, tt.expectedResponse, recorder.Body.String())
        })
    }
}
