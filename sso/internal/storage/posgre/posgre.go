package posgre

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

var (
	Table = "users"
	Apps  = "apps"
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.posgre.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s.Ping: %w", op, err)
	}

	// возможно здесь запустить миграции как в п.посгре

	return &Storage{db: db}, nil
}

// CloseDB close database
func (storage *Storage) CloseDB() {
	// storage.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.posgre.SaveUser"

	var id int64

	query := fmt.Sprintf("INSERT INTO  %s  (email, pass_hash) VALUES ($1,$2) RETURNING id", Table)
	err := s.db.QueryRow(
		query, email, passHash,
	).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		// Проверка на уникальность
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return 0, fmt.Errorf("U:%s:  %w", op, err)
		}
		return 0, fmt.Errorf("O:%s:  %w", op, err)

	}

	return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {

	const op = "storage.posgre.User"

	var user models.User
	query := fmt.Sprintf("SELECT id, email, pass_hash FROM %s WHERE email = $1", Table)
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("N:%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("O:%s:  %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.posgre.IsAdmin"

	var isAdmin bool

	query := fmt.Sprintf("SELECT is_admin FROM  FROM %s WHERE is_admin = $1", Table)
	err := s.db.QueryRow(query, userID).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s:  %w", op, err)
	}

	return isAdmin, nil
}

// func (s *Storage) CLose()  error  {
// 	reurn
// }

func (s *Storage) App(ctx context.Context, id int64) (models.App, error) {
	const op = "storage.posgre.App"

	var app models.App
	query := fmt.Sprintf("SELECT id, name, secret FROM %s WHERE id = $1", Apps)
	err := s.db.QueryRow(query, id).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("N:%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("O:%s:  %w", op, err)
	}
	return app, nil
}
