package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/yolkhovyy/go-otelw/pkg/slogw"
	"github.com/yolkhovyy/go-userv/internal/contract/storage"
)

func (c *Controller) Create(ctx context.Context, user storage.UserInput) (*storage.User, error) {
	query := `
		INSERT INTO users (id, first_name, last_name, nickname, password_hash, email, country, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, first_name, last_name, nickname, email, country, created_at, updated_at`

	row := c.pool.QueryRow(ctx, query,
		uuid.New(), user.FirstName, user.LastName, user.Nickname,
		user.Password, user.Email, user.Country)

	createdUser := &storage.User{}
	if err := row.Scan(&createdUser.ID, &createdUser.FirstName, &createdUser.LastName,
		&createdUser.Nickname, &createdUser.Email, &createdUser.Country,
		&createdUser.CreatedAt, &createdUser.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

func (c *Controller) Update(ctx context.Context, user storage.UserUpdate) (*storage.User, error) {
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

	row := c.pool.QueryRow(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Nickname,
		user.Password, user.Email, user.Country)

	updatedUser := storage.User{}
	if err := row.Scan(&updatedUser.ID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Nickname,
		&updatedUser.Email, &updatedUser.Country, &updatedUser.CreatedAt, &updatedUser.UpdatedAt); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &updatedUser, nil
}

func (c *Controller) Get(ctx context.Context, userID uuid.UUID) (*storage.User, error) {
	var user storage.User

	query := `SELECT id, first_name, last_name, nickname, email, country, created_at, updated_at
		FROM users WHERE id = $1`

	if err := pgxscan.Get(ctx, c.pool, &user, query, userID); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil
}

func (c *Controller) List(ctx context.Context, page int, limit int, countryCode string) ([]storage.User, int, error) {
	logger := slogw.DefaultLogger()

	trx, err := c.txBeginRepeatableRead(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list users transaction: %w", err)
	}

	defer func() {
		if err := trx.Rollback(ctx); err != nil {
			logger.ErrorContext(ctx, "transaction",
				slog.String("rollback", err.Error()),
			)
		}
	}()

	var count int

	query := `SELECT COUNT(*) FROM users WHERE country = $1`

	if err := pgxscan.Get(ctx, trx, &count, query, countryCode); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	var args []any

	offset := (page - 1) * limit
	query = `SELECT id, first_name, last_name, nickname, email, country, created_at, updated_at FROM users `

	if countryCode != "" {
		args = append(args, countryCode, limit, offset)
		query += `WHERE country = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	} else {
		args = append(args, limit, offset)
		query += `ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	var result []storage.User
	if err := pgxscan.Select(ctx, trx, &result, query, args...); err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	return result, count, nil
}

func (c *Controller) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := c.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

//nolint:ireturn
func (c *Controller) txBeginRepeatableRead(ctx context.Context) (pgx.Tx, error) {
	trx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	_, err = trx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		if rbErr := trx.Rollback(ctx); rbErr != nil {
			return nil, fmt.Errorf("set isolation level: %w, rollback transaction: %w", err, rbErr)
		}

		return nil, fmt.Errorf("set isolation level: %w", err)
	}

	return trx, nil
}
