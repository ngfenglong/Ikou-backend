package repository

import (
	"context"
	"time"

	"github.com/ngfenglong/ikou-backend/api/dto"
	"github.com/ngfenglong/ikou-backend/api/models"
)

func (m *DBModel) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u models.User
	row := m.DB.QueryRowContext(ctx, `
		SELECT id, username, email, password, firstname, lastname, country, profileImage
		FROM Users
		WHERE username = ?
	`, username)

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.FirstName,
		&u.LastName,
		&u.Country,
		&u.ProfileImage,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
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

func (m *DBModel) RegisterUser(r dto.RegisterFormInputDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		Insert into Users (username, firstname, lastname, email, country, profileImage, password) 
		VALUE (?,?,?,?,?,?,?)
	`
	_, err := m.DB.ExecContext(ctx, stmt, r.Username, r.FirstName, r.LastName, r.Email, r.Country, r.ProfileImage, r.Password)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) CheckIfUserExists(r dto.RegisterFormInputDTO) (bool, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	checkUsernameStmt := `
		Select 
			Count(*) 
		From 
			Users 
		WHERE
			username = ?
	`

	row := m.DB.QueryRowContext(ctx, checkUsernameStmt, r.Username)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, false, err
	}
	if count > 0 {
		return true, false, nil
	}

	checkEmailStmt := `
		Select 
			Count(*) 
		From 
			Users 
		WHERE
			email = ?
	`
	row = m.DB.QueryRowContext(ctx, checkEmailStmt, r.Email)

	err = row.Scan(&count)
	if err != nil {
		return false, false, err
	}
	if count > 0 {
		return false, true, nil
	}

	return false, false, nil
}

func (m *DBModel) DeleteToken(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `Delete from RefreshTokens where token = ?`
	_, err := m.DB.ExecContext(ctx, stmt, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
