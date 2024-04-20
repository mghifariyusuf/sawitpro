// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/handler/models"
	"github.com/SawitProRecruitment/UserService/repository/entity"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, req *models.RegisterUserRequest) (int64, error)
	GetUser(ctx context.Context, filter *entity.UserFilter) (*entity.UserData, error)
	UpdateProfile(ctx context.Context, userID int64, req *models.UpdateUserProfileRequest) error
	IncLogin(ctx context.Context, userID int64) error
}
