package user

import (
	"context"
	"database/sql"
	"ddd2/exception"
	"ddd2/layer/auth"
	"ddd2/model/domain"
	"ddd2/model/web"
	"ddd2/utils"
)

type ServiceImpl struct {
	Repo Repository
	Db   *sql.DB
}

func NewService(repo Repository, db *sql.DB) Service {
	return &ServiceImpl{
		Repo: repo,
		Db:   db,
	}
}

type Service interface {
	GetByUsername(ctx context.Context, username string) (*web.User, error)
	GetByEmail(ctx context.Context, email string) (*web.User, error)
	GetById(ctx context.Context, id int) (*web.User, error)
	CreateUser(ctx context.Context, user web.UserRegisterPayload) error
}

func (s *ServiceImpl) GetByUsername(ctx context.Context, username string) (*web.User, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	u, err := s.Repo.GetByUsername(ctx, tx, username)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserToWeb(u), nil
}

func (s *ServiceImpl) GetByEmail(ctx context.Context, email string) (*web.User, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	u, err := s.Repo.GetByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserToWeb(u), nil
}

func (s *ServiceImpl) GetById(ctx context.Context, id int) (*web.User, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	u, err := s.Repo.GetById(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserToWeb(u), nil
}

func (s *ServiceImpl) CreateUser(ctx context.Context, user web.UserRegisterPayload) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	_, err = s.GetByUsername(ctx, user.Username)
	if err == nil {
		return exception.ErrUsernameTaken
	}

	_, err = s.GetByEmail(ctx, user.Email)
	if err == nil {
		return exception.ErrEmailTaken
	}

	p, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	User := &domain.User{
		Username: user.Username,
		Password: p,
		Email:    user.Email,
		Name:     user.Name,
	}

	err = s.Repo.CreateUser(ctx, tx, *User)
	if err != nil {
		return err
	}

	return nil
}
