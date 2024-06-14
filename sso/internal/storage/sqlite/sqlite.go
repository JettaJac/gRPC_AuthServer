package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go/constant"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "qithub.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of SqLite Storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage(db:db), nil
}

func (s *Storage) SaveUser(ctx.Context, email string, passHash []byte) (int64, error)  {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?,?);")
	if err != nil  {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx,email, passHash)
	if err != nil  {
		var sqlite sqlite3.Error
		// Проверка на уникальность
		 if errors.As(err, &sqlite.Error) && sqliteErr.ExtendedCode == sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique{
			return 0, fmt.Errorf("%s:  %w", op, err)
		 }
return 0, fmt.Errorf("%s:  %w", op, err)

	}

// Получаем  ID созданного пользователя
id, err  := res.LastInsertId()
if err != nil   {
	return 0, fmt.Errorf("%s:  %w", op, err)
}
return id, nil
}

//User returns user by email
func (s *Storage) User(ctx context.Context, email string)  (model.User, error){

const op  = "storage.sqlite.User"

stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM user WHERE email =?")
	if err != nil  {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	} 

	row := stmt.QueryRowContext(ctx, email)

	var user model.User
	err= row.Scan(&user.ID,&user.Email,&user.PassHash)
	if err != nil  {
		if errors.Is(err,  sql.ErrNoRows)  {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)		
		}
		return model.User{}, fmt.Errorf("%s:  %w", op, err)
	}

return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error)  {
	const op  =  "storage.sqlite.IsAdmin"
	stmt, err := s.db.Prepare("SELECT is_admin FROM user WHERE is_admin =?")
	if err != nil  {
		return false, fmt.Errorf("%s: %w", op, err)
	} 
	row := stmt.QueryRowContext(ctx, userID)
	var isAdmin  bool
	err= row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err,  sql.ErrNoRows)  {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)		
		}
		return false, fmt.Errorf("%s:  %w", op, err)	
	}

	return isAdmin, nil
}

// func (s *Storage) CLose()  error  {
// 	reurn
// }

func (s *Storage) App(ctx context.Context, appID int)  (model.App, error)   {
	const op =  "storage.sqlite.App"

	stmt, err  := s.db.Prepare("SELECT id, name, secret FROM app WHERE app_id  =?")
	if err != nil   {
		return model.App{}, fmt.Errorf("%s: %w", op, err)	
	}

	row  := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&appID, &app.Name,&app.Secret)
	if err != nil   {
		if errors.Is(err,  sql.ErrNoRows)  {
			return model.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)		
		}
		return model.App{}, fmt.Errorf("%s:  %w", op, err)	
	}
}