package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	sl "sso/internal/lib/logger"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, uid int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, uid int64) (models.App, error)
}

/// Как раз тут про кеш 1.42

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new Auth instance.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProviderchan AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProviderchan,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if  user with given credentials exists in the system.
//
// If user exists, but password is wrong, returns an error.
// If user doesn't exist, returns an error.
func (a *Auth) Login(
	ctx context.Context,
	email,
	password string,
	appID int,
) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email), //!!! не логирровать персональные данные
	)

	// check user
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s:  %v", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s:  %v", op, err)
	}

	// check password
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s:  %v", op, ErrInvalidCredentials)
	}

	// check app
	app, err := a.appProvider.App(ctx, appID) //лучше хранить ключи не в бд
	if err != nil {
		return "", fmt.Errorf("%s:  %v", op, err)
	}

	log.Info("user logger in successfuly")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("%s:  %v", op, err)
	}

	return token, nil
}

// RegisterNewUser creates a new user in the system and returns its ID.
// If user with given username already exists, returns an error.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("email", email), // так хранитьб логировать не стоить именно персональные данные
	)
	log.Info("registring user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s:  %v", op, ErrUserExists)
		}
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return false, fmt.Errorf("%s:  %v", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s:  %v", op, err)
	}

	log.Info("user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}
