package user

import (
	"database/sql"

	"github.com/EmiliodDev/todoAPI/types"
)

type Store struct {
    db *sql.DB
}

func NewStore(db *sql.DB) *Store {
    return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
    return nil, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
    return nil, nil
}

func (s *Store) CreateUser(user types.User) error {
    return nil
}
