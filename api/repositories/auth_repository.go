package repository

import (
	"context"
	"time"

	"github.com/ngfenglong/ikou-backend/api/models"
)

func (m *DBModel) GetUserByUsername(username string) (models.LoginUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.LoginUser
	row := m.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password
		FROM Users
		WHERE username = ?
	`, username)

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *DBModel) InsertToken(userId string, refreshToken string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// delete existing tokens
	stmt := `Delete from RefreshTokens where userId = ?`
	_, err := m.DB.ExecContext(ctx, stmt, userId)
	if err != nil {
		return err
	}

	stmt = `INSERT into RefreshTokens (userId, token, expires_at, created_at, updated_at) 
		values (?, ?, ?, ?, ?)
	`

	_, err = m.DB.ExecContext(ctx, stmt, userId, refreshToken, expiresAt, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil

}
