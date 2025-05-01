package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
	"github.com/Fadil-Tao/gopher-pay/utils/httpresponse"
	"github.com/golang-jwt/jwt/v5"
)

type key string

const tokenKey key = "user"

func VerifyToken(tokenString string) (*jwt.Token,error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		slog.Error("error parsing jwt token", "message" , err)
		return nil,csterr.ErrInternal
	}

	if !token.Valid {
		slog.Error("jwt token is not valid", "message", err)
		return nil,csterr.ErrUnauthorized
	}
	return token, nil
}


func RequireAuth(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			httpresponse.WriteErrorResponse(w,"unauthorized",nil,"it's seems you're not logged in yet",http.StatusUnauthorized)			
			return 
		}
		tokenValue := cookie.Value
		token, err := VerifyToken(tokenValue)
		if err != nil {
			httpresponse.WriteErrorResponse(w,"unauthorized",nil,"unautorizhed",http.StatusUnauthorized)
			return 
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := r.Context()

			ctx = context.WithValue(ctx, tokenKey, claims)
			r = r.WithContext(ctx)
		} else {
			slog.Error("error parsing token claims", "message", err)
			httpresponse.WriteErrorResponse(w,"unauthorized",nil,"invalid token",http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	claims, ok := r.Context().Value(tokenKey).(jwt.MapClaims)
	if !(ok) {
		return 0, csterr.ErrUnauthorized
	}

	userIdFloat, ok := claims["id"].(float64)
	if !ok {
		return 0,csterr.ErrUnauthorized
	}
	userId := int(userIdFloat)
	return userId, nil
}