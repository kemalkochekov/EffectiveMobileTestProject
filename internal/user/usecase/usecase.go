package usecase

import (
	"GraphQL/configs"
	"GraphQL/internal/models"
	"GraphQL/internal/user"
	"GraphQL/pkg/errlst"
	"GraphQL/pkg/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strconv"
)

type UserUC struct {
	cfg    *configs.Config
	pgRepo user.Repository
}

func NewUserUC(cfg *configs.Config, pgDB user.Repository) *UserUC {
	return &UserUC{
		cfg:    cfg,
		pgRepo: pgDB,
	}
}

func (u *UserUC) GetAgeDataByName(_ context.Context, name string) (int, error) {
	resp, err := http.Get("https://api.agify.io/?name=" + name)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var agifyResponse models.AgifyResponse
	agifyResponse.Age = -1

	err = json.NewDecoder(resp.Body).Decode(&agifyResponse)
	if err != nil {
		return -1, err
	}

	if agifyResponse.Age == -1 {
		return -1, errlst.ErrEmptyAgify
	}

	return agifyResponse.Age, nil
}
func (u *UserUC) GetGenderDataByName(_ context.Context, name string) (string, error) {
	resp, err := http.Get("https://api.genderize.io/?name=" + name)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var genderizeResponse models.GenderizeResponse

	err = json.NewDecoder(resp.Body).Decode(&genderizeResponse)
	if err != nil {
		return "", err
	}

	if genderizeResponse.Gender == "" {
		return "", errlst.ErrEmptyGenderize
	}

	return genderizeResponse.Gender, nil
}
func (u *UserUC) GetNationalityDataByName(_ context.Context, name string) (string, error) {
	resp, err := http.Get("https://api.nationalize.io/?name=" + name)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var nationalizeResponse models.NationalizeResponse

	err = json.NewDecoder(resp.Body).Decode(&nationalizeResponse)
	if err != nil {
		return "", err
	}

	if nationalizeResponse.Country == nil {
		return "", errlst.ErrEmptyNationalize
	}
	countryID := ""
	minVal := -1.0

	for _, val := range nationalizeResponse.Country {
		// The big package supports big numbers and it supports Int, Rational Number and Float.
		// or we can use tolerance := 0.0001 both methods is appropriate
		result := big.NewFloat(minVal).Cmp(big.NewFloat(val.Probability))
		if result < 0 {
			minVal = val.Probability
			countryID = val.CountryID
		}
	}

	return countryID, nil
}

func (u *UserUC) CreateIndividualData(
	ctx context.Context,
	ageResp int,
	genderResp string,
	nationalizeResp string,
	name string,
	surname string,
	patronymic *string,
) (models.Individual, error) {
	var individual models.Individual
	individual.Age = &ageResp
	individual.Gender = &genderResp
	individual.CountryID = &nationalizeResp
	individual.Name = name
	individual.Surname = surname

	if patronymic != nil {
		individual.Patronymic = patronymic
	}

	checkID, err := u.pgRepo.CheckUserExists(ctx, name, surname, patronymic)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return models.Individual{}, err
	}

	if checkID != -1 {
		return models.Individual{}, errlst.ErrUserAlreadyExists
	}

	insertedID, err := u.pgRepo.CreateUserDataByName(ctx, individual)
	if err != nil {
		return models.Individual{}, err
	}

	individual.ID = insertedID

	return individual, nil
}

func (u *UserUC) UpdateAllUserDataByID(
	ctx context.Context,
	idNumber int64,
	age int,
	gender string,
	countryID string,
	name *string,
	surname *string,
	patronymic *string,
) error {
	var individual models.Individual

	userData, err := u.pgRepo.GetUserDataByID(ctx, idNumber)
	if err != nil {
		return err
	}

	// Checking Surname changed
	individual.Surname = userData.Surname
	if surname != nil && userData.Surname != *surname && *surname != "" {
		individual.Surname = *surname
	}

	// Checking Patronymic changed
	individual.Patronymic = userData.Patronymic
	if patronymic != nil && userData.Patronymic != patronymic {
		individual.Patronymic = patronymic
	}

	individual.ID = strconv.Itoa(int(idNumber))
	individual.Age = &age
	individual.Gender = &gender
	individual.CountryID = &countryID
	individual.Name = userData.Name
	// Checking name is changed otherwise we need to update all information of user
	if userData.Name != *name {
		individual.Name = *name

		err = u.pgRepo.UpdateUserAllDataByID(ctx, individual)
		if err != nil {
			return err
		}

		return nil
	}

	err = u.pgRepo.UpdateUserDataByID(ctx, idNumber, individual.Surname, individual.Patronymic)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUC) UpdateIndividualDataByID(
	ctx context.Context,
	idNumber int64,
	surname *string,
	patronymic *string,
) error {
	userData, err := u.pgRepo.GetUserDataByID(ctx, idNumber)
	if err != nil {
		return err
	}

	// Checking Surname changed
	if surname != nil && userData.Surname != *surname && *surname != "" {
		userData.Surname = *surname
	}

	// Checking Patronymic changed
	if patronymic != nil && userData.Patronymic != patronymic {
		userData.Patronymic = patronymic
	}

	err = u.pgRepo.UpdateUserDataByID(ctx, idNumber, userData.Surname, userData.Patronymic)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUC) DeleteAllUserDataByID(ctx context.Context, idNumber int64) error {
	err := u.pgRepo.DeleteAllUserDataByID(ctx, idNumber)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUC) GetUserData(ctx context.Context, idNumber int64) (*models.Individual, error) {
	resp, err := u.pgRepo.GetUserAllDataByID(ctx, idNumber)
	if err != nil {
		return nil, err
	}

	return resp.ToServer(), nil
}

func (u *UserUC) GetUserDataByFilter(
	ctx context.Context,
	filter *models.IndividualFilter,
	page *int,
	pageSize *int,
	sortBy *string,
) ([]*models.Individual, error) {
	resp, err := u.pgRepo.GetUserDataByFilter(ctx, filter, page, pageSize, sortBy)
	if err != nil {
		return nil, err
	}

	responseData := utils.Map(
		resp,
		func(item *models.IndividualData) *models.Individual {
			return item.ToServer()
		},
	)

	return responseData, nil
}
