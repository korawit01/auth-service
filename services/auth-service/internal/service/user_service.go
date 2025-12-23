package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/korawit01/auth-service/services/auth-service/internal/domain"
	"github.com/korawit01/auth-service/services/auth-service/internal/repository"
	"github.com/korawit01/auth-service/services/auth-service/internal/token"
)

type UserService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error) // return JWT
	VerifyToken(ctx context.Context, tokenStr string) (*token.Claims, error)
}

type userService struct {
	repo      repository.UserRepository
	tokenProv token.Provider
}

func NewUserService(r repository.UserRepository, t token.Provider) UserService {
	return &userService{repo: r, tokenProv: t}
}

func (s *userService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Email:    email,
		Password: string(hash),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", err
	}

	// build JWT token for the authenticated user
	return s.tokenProv.Generate(u.RowId, u.Email)
}

func (s *userService) VerifyToken(ctx context.Context, tokenStr string) (*token.Claims, error) {
	_ = ctx
	return s.tokenProv.Verify(tokenStr)
}
