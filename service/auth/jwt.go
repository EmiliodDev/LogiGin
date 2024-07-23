package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/EmiliodDev/LogiGin/config"
	"github.com/EmiliodDev/LogiGin/types"
	"github.com/EmiliodDev/LogiGin/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handlerFunc gin.HandlerFunc, store types.UserStore) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := utils.GetTokenFromRequest(c)

        token, err := validateJWT(tokenString)
        if err != nil {
            log.Printf("failed to validate token: %v", err)
            permissionDenied(c)
            return
        }

        if !token.Valid {
            log.Println("invalid token")
            permissionDenied(c)
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        str := claims["userID"].(string)

        userID, err := strconv.Atoi(str)
        if err != nil {
            log.Printf("failed to convert userID to int: %v", err)
            permissionDenied(c)
            return
        }

        u, err := store.GetUserByID(userID)
        if err != nil {
            log.Printf("failed to get user by id: %v", err)
            permissionDenied(c)
            return
        }
        
        ctx := context.WithValue(c.Request.Context(), UserKey, u.ID)
        c.Request = c.Request.WithContext(ctx)

        handlerFunc(c)
    }
}

func CreateJWT(secret []byte, userID int) (string, error) {
    expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": strconv.Itoa(int(userID)),
        "expiresAt": time.Now().Add(expiration).Unix(),
    })

    tokenString, err := token.SignedString(secret)
    if err != nil {
        return "", err
    }

    return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        return []byte(config.Envs.JWTSecret), nil
    })
}

func permissionDenied(c *gin.Context) {
    c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
}
