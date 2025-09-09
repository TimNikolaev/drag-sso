package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/TimNikolaev/drag-sso/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserCreator interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (uint64, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int32) (*models.App, error)
}

type Auth struct {
	log          *slog.Logger
	tokenTTL     time.Duration
	userCreator  UserCreator
	userProvider UserProvider
	appProvider  AppProvider
}

func New(
	log *slog.Logger,
	tokenTTL time.Duration,
	userCreator UserCreator,
	userProvider UserProvider,
	appProvider AppProvider,
) *Auth {
	return &Auth{
		log:          log,
		tokenTTL:     tokenTTL,
		userCreator:  userCreator,
		userProvider: userProvider,
		appProvider:  appProvider,
	}
}

//Realization Auth-service's methods

func (a *Auth) SignUpNewUser(ctx context.Context, email string, pass string) (uint64, error) {
	const op = "auth.SignUpNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("sign upping user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hashing password")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userCreator.CreateUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to create user")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created user")

	return id, nil

}

func (a *Auth) SignIn(ctx context.Context, email string, pass string) (string, error) {
	panic("not implemented")
}
