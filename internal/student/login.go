package student

import (
	"context"
	"errors"
	utils "golang-assignment/utils"
)

type Services struct {
	Store StudentStore
}

type User struct {
	ID       string
	Password string
}

type StudentService interface {
	AuthenticateUser(ctx context.Context, userID, password string) (User, error)
	GenerateJWT(user User) (string, error)
}

func (s *Service) AuthenticateUser(ctx context.Context, userID, password string) (User, error) {
	if userID == "user123" && password == "password" {
		return User{ID: "user123"}, nil
	}
	return User{}, errors.New("invalid credentials")
}

func (s *Service) GenerateJWT(user User) (string, error) {
	return utils.GenerateJWT(user.ID)
}
