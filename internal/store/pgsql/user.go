package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"mistar-be-go/internal/store"
	"time"

	"github.com/google/uuid"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

const userInsert = `INSERT INTO
users(
	id, type, name, email, phone_number, password, created_at
) values(
	$1, $2, $3, $4, $5, $6, $7	
)
`

func (s *User) Insert(ctx context.Context, user *store.UserData) error {
	insertStmt, err := s.db.PrepareContext(ctx, userInsert)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC()
	userID := uuid.NewString()
	_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
		userID, user.Type, user.Name,
		user.Email, user.PhoneNumber, user.Password, createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	user.ID = userID

	return nil
}

const FindUserByEmailQuery = `SELECT u.id, u.type, u.name, u.email, u.phone_number, u.password
	FROM users u
	WHERE u.email = $1
	LIMIT 1
`

func (s *User) FindUserByEmail(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, FindUserByEmailQuery, email)

	err := row.Scan(
		&user.ID, &user.Type, &user.Name, &user.Email, &user.PhoneNumber, &user.Password,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
