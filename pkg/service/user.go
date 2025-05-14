package service

import (
	"log/slog"

	"megumin/pkg/lib/sql"
	"megumin/pkg/service/internal/repository"
)

type User interface {
	Create(id string) error
	Exists(id string) bool
}

type user struct {
	userRepository repository.User
	log            *slog.Logger
}

func NewUser(db sql.Connection, log *slog.Logger) User {
	return &user{
		userRepository: repository.NewUser(db),
		log:            log,
	}
}

func (u user) Create(id string) error {
	return u.userRepository.Create(id)
}

func (u user) Exists(id string) bool {
	exists, err := u.userRepository.Exists(id)
	if err != nil {
		u.log.Error("Error checking if user exists", "error", err)
		return false
	}

	return exists
}
