package user

import (
	"GraphQL/internal/models"
	"context"
)

type UseCase interface {
	GetAgeDataByName(ctx context.Context, name string) (int, error)
	GetGenderDataByName(ctx context.Context, name string) (string, error)
	GetNationalityDataByName(ctx context.Context, name string) (string, error)
	CreateIndividualData(
		ctx context.Context,
		ageResp int,
		genderResp string,
		nationalizeResp string,
		name,
		surname string,
		patronymic *string,
	) (models.Individual, error)

	UpdateAllUserDataByID(
		ctx context.Context,
		idNumber int64,
		ageResp int,
		genderResp string,
		nationalityResp string,
		name *string,
		surname *string,
		patronymic *string,
	) error

	UpdateIndividualDataByID(ctx context.Context, idNumber int64, surname *string, patronymic *string) error
	DeleteAllUserDataByID(ctx context.Context, idNumber int64) error
	GetUserData(ctx context.Context, idNumber int64) (*models.Individual, error)
	GetUserDataByFilter(
		ctx context.Context,
		filter *models.IndividualFilter,
		page *int,
		pageSize *int,
		sortBy *string,
	) ([]*models.Individual, error)
}
