package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/TimNikolaev/drag-sso/internal/models"
	"github.com/TimNikolaev/drag-sso/internal/repository"
	"github.com/TimNikolaev/drag-sso/pkg/jwt"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *Auth) SignUpNewUser(ctx context.Context, email string, pass string) (uint64, error) {
	const op = "auth.SignUpNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("sign upping user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hashing password" /**/)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userCreator.CreateUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to create user" /**/)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("created user")

	return id, nil

}

func (a *Auth) SignIn(ctx context.Context, email string, pass string, appID int32) (string, error) {
	const op = "auth.SignIn"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attempting to sign in user")

	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Error("user not found" /**/)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user" /**/)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		log.Error("invalid credentials" /**/)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user sign in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token" /**/)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
