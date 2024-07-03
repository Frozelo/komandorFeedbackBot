package service

import (
	"sync"

	"github.com/Frozelo/komandorFeedbackBot/internal/domain/entity"
)

type UserRepository interface {
	CreateTgUser(user entity.User) (entity.User, error)
}

type UserService struct {
	repo UserRepository
	mu   sync.Mutex
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) CreateUser(user entity.User) (entity.User, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	return us.repo.CreateTgUser(user)
}
