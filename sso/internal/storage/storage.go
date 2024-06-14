package storage

import (
	"errors"
)

// migrate create  -ext sql -dir migrations --seq create_commands  - команда файлы с миграциями
var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)
