package postgres

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/yolkhovyy/user/internal/contract/storage"
)

func (u *Controller) Create(ctx context.Context, user storage.UserInput) (*storage.User, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `
		INSERT INTO users (id, first_name, last_name, nickname, password_hash, email, country, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, first_name, last_name, nickname, email, country, created_at, updated_at`

	row := u.pool.QueryRow(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Nickname,
		user.Password, user.Email, user.Country)

	createdUser := &storage.User{}
	if err := row.Scan(&createdUser.ID, &createdUser.FirstName, &createdUser.LastName,
		&createdUser.Nickname, &createdUser.Email, &createdUser.Country,
		&createdUser.CreatedAt, &createdUser.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

func (u *Controller) Update(ctx context.Context, user storage.UserInput) (*storage.User, error) {
	query := `
		UPDATE users SET
			first_name = CASE WHEN $2 != '' THEN $2 ELSE users.first_name END,
			last_name = CASE WHEN $3 != '' THEN $3 ELSE users.last_name END,
			nickname = CASE WHEN $4 != '' THEN $4 ELSE users.nickname END,
			password_hash = CASE WHEN $5 != '' THEN $5 ELSE users.password_hash END,
			email = CASE WHEN $6 != '' THEN $6 ELSE users.email END,
			country = CASE WHEN $7 != '' THEN $7 ELSE users.country END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, first_name, last_name, nickname, email, country, created_at, updated_at`

	row := u.pool.QueryRow(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Nickname,
		user.Password, user.Email, user.Country)

	updatedUser := &storage.User{}
	if err := row.Scan(&user.ID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Nickname,
		&updatedUser.Email, &updatedUser.Country, &updatedUser.CreatedAt, &updatedUser.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return updatedUser, nil
}

func (u *Controller) Get(ctx context.Context, userID uuid.UUID) (*storage.User, error) {
	var user storage.User

	query := `SELECT id, first_name, last_name, nickname, email, country, created_at, updated_at
		FROM users WHERE id = $1`

	if err := pgxscan.Get(ctx, u.pool, &user, query, userID); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil
}

func (u *Controller) List(ctx context.Context, page int, limit int, countryCode string) ([]storage.User, error) {
	var args []any

	offset := (page - 1) * limit
	query := `SELECT id, first_name, last_name, nickname, email, country, created_at, updated_at FROM users `

	if countryCode != "" {
		args = append(args, countryCode, limit, offset)
		query += `WHERE country = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	} else {
		args = append(args, limit, offset)
		query += `ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	var result []storage.User
	if err := pgxscan.Select(ctx, u.pool, &result, query, args...); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return result, nil
}

func (u *Controller) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := u.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (u *Controller) Count(ctx context.Context) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM users`

	if err := pgxscan.Get(ctx, u.pool, &count, query); err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}
