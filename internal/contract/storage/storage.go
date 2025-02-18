package storage

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type Contract interface {
	CRUD
	io.Closer
}

type CRUD interface {
	Creator
	Reader
	Updater
	Deleter
}

type Creator interface {
	Create(ctx context.Context, user UserInput) (*User, error)
}

type Reader interface {
	Get(ctx context.Context, userID uuid.UUID) (*User, error)
	List(ctx context.Context, page, limit int, countryCode string) ([]User, int, error)
}

type Updater interface {
	Update(ctx context.Context, user UserUpdate) (*User, error)
}

type Deleter interface {
	Delete(ctx context.Context, userID uuid.UUID) error
}

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Email     string    `json:"email" db:"email"`
	Country   string    `json:"country" db:"country"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Extra field: Password.
type UserInput struct {
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Nickname  string `json:"nickname" db:"nickname"`
	Email     string `json:"email" db:"email"`
	Country   string `json:"country" db:"country"`
	Password  string `json:"password,omitempty" db:"password_hash"`
}

type UserUpdate struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Email     string    `json:"email" db:"email"`
	Country   string    `json:"country" db:"country"`
	Password  string    `json:"password,omitempty" db:"password_hash"`
}
type UserList struct {
	Users      []User `json:"users,omitempty"`
	TotalCount int    `json:"totalCount"`
	NextPage   int    `json:"nextPage"`
}
