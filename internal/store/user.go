package store

import "context"

type UserData struct {
	ID          string
	Type        string
	Name        string
	Email       string
	Password    string
	PhoneNumber string
}

type User interface {
	Insert(ctx context.Context, user *UserData) error
	FindUserByEmail(ctx context.Context, email string) (*UserData, error)
}
