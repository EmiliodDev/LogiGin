package types

import "time"

type UserStore interface {
    GetUserByEmail(email string) (*User, error)
    GetUserByID(id int) (*User, error)
    CreateUser(User) error
}

type User struct {
    ID          int         `json:"id"`
    FirstName   string      `json:"firstName"`
    LastName    string      `json:"lastName"`
    Email       string      `json:"email"`
    Password    string      `json:"password"`
    CreatedAt   time.Time   `json:"createdAt"`
}

type RegisterUserPayload struct {
    FirstName   string      `json:"firstName" validate:"required"`
    LastName    string      `json:"lastName" validate:"required"`
    Email       string      `json:"email" validate:"required,email"`
    Password    string      `json:"password" validate:"required,min=3,max=30"`
}

type LoginUserPayload struct {
    Email       string      `json:"email" validate:"required,email"`
    Password    string      `json:"password" validate:"required"`
}
