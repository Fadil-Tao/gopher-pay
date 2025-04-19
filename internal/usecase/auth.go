package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Fadil-Tao/gopher-pay/internal/model"
	csterr "github.com/Fadil-Tao/gopher-pay/utils/errors"
	"github.com/Fadil-Tao/gopher-pay/utils/hashing"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRepo interface {
	Register(ctx context.Context, user model.User) error
	IsUserExist(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

var argon2idHash = hashing.NewArgon2idHash(1, 16, 64*1024, 1, 32)

type AuthService struct {
	AuthRepo
}

func NewAuthSerice(authRepo AuthRepo) *AuthService {
	return &AuthService{
		authRepo,
	}
}

func (a *AuthService) Register(ctx context.Context, user model.User) error {
	select {
	case <-ctx.Done():
		slog.Error("context passed")
		return ctx.Err()
	default:
	}

	isUserExist, err := a.AuthRepo.IsUserExist(ctx, user.Email)
	if err != nil {
		slog.Error(err.Error())
		return csterr.ErrInternal
	}
	if isUserExist {
		return csterr.ErrIsAlreadyExist
	}

	hashedPassword, err := argon2idHash.GenerateHash(user.Password, nil)
	if err != nil {
		slog.Error(err.Error())
		return csterr.ErrInternal
	}

	encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword.Hash)
	encodedSalt := base64.StdEncoding.EncodeToString(hashedPassword.Salt)

	user.Password = encodedPassword
	user.Salt = encodedSalt

	err = a.AuthRepo.Register(ctx, user)
	if err != nil {
		return csterr.ErrInternal
	}
	return nil
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (response *string, err error) {
	user, err := a.AuthRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		slog.Error("error decoding password", "message", err)
		err = csterr.ErrInternal
		return nil, err
	}

	decodedSalt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		slog.Error("error decoding salt", "message", err)
		err = csterr.ErrInternal
		return nil, err
	}

	isMatch, err := argon2idHash.Compare(decodedPassword, decodedSalt, user.Password)
	if err != nil {
		slog.Error("error comparing password", "message", err)
		return nil, nil
	}
	if !isMatch {
		return nil, fmt.Errorf("invalid credentials")
	}

	jwtToken, err := createToken(user)
	if err != nil {
		slog.Error("failed to create jwt token", "message", err)
		err = csterr.ErrInternal
		return nil, err
	}
	return &jwtToken, nil
}

func createToken(user *model.User) (string, error) {
	expiry := time.Now().Add(time.Hour * 168).Unix()
	secretKey := os.Getenv("JWT_SECRET")
	key := []byte(secretKey)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
			"exp":   expiry,
		})
	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}
