package user

import (
	"GraphQL/internal/models"
	"context"
)

type Repository interface {
	CreateUserDataByName(ctx context.Context, individual models.Individual) (string, error)
	CheckUserExists(ctx context.Context, name, surname string, patronymic *string) (int64, error)
	GetUserDataByID(ctx context.Context, idNumber int64) (models.UserData, error)
	DeleteAllUserDataByID(ctx context.Context, idNumber int64) error
	UpdateUserAllDataByID(ctx context.Context, individual models.Individual) error
	UpdateUserDataByID(ctx context.Context, idNumber int64, surname string, patronymic *string) error
	GetUserAllDataByID(ctx context.Context, idNumber int64) (*models.IndividualData, error)
	GetUserDataByFilter(
		ctx context.Context,
		filter *models.IndividualFilter,
		page *int,
		pageSize *int,
		sortBy *string,
	) ([]*models.IndividualData, error)
}
