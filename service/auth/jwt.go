package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/m21power/Ecom/config"
	"github.com/m21power/Ecom/types"
	"github.com/m21power/Ecom/utils"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, UserId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(UserId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the user
		tokenString := getTokenFromRequest(r)
		if tokenString == "" {
			log.Println("No token found in request")
			permissionDenied(w)
			return
		}

		// Validate the JWT
		token, err := validateToken(tokenString)
		if err != nil {
			log.Println("Error validating token:", err)
			permissionDenied(w)
			return
		}
		if !token.Valid {
			log.Println("Token is not valid")
			permissionDenied(w)
			return
		}

		// If it is valid, fetch user ID from the DB
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		userId, _ := strconv.Atoi(str)
		u, err := store.GetUserByID(userId)
		if err != nil {
			log.Println("Error fetching user by ID:", err)
			permissionDenied(w)
			return
		}

		// Set context with the user ID
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	// Get the token from the request
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth == "" {
		return ""
	}

	// Ensure the token is prefixed with "Bearer "
	if !strings.HasPrefix(tokenAuth, "Bearer ") {
		return ""
	}

	return strings.TrimPrefix(tokenAuth, "Bearer ")
}

func validateToken(t string) (*jwt.Token, error) {
	// Check if the token has three parts
	parts := strings.Split(t, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("token is malformed: expected 3 parts, got %d", len(parts))
	}

	// Validate the token
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetIdFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}
