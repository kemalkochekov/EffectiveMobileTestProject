package postgres

import (
	"GraphQL/configs"
	"GraphQL/internal/models"
	"GraphQL/pkg/connection"
	"GraphQL/pkg/errlst"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type UserRepo struct {
	cfg *configs.Config
	db  connection.DBops
}

func NewClientPGRepository(
	cfg *configs.Config,
	db connection.DBops,
) *UserRepo {
	return &UserRepo{cfg: cfg, db: db}
}
func (u *UserRepo) CheckUserExists(ctx context.Context, name, surname string, patronymic *string) (int64, error) {
	query := "SELECT id FROM individuals WHERE name = $1 AND surname = $2"

	if patronymic != nil {
		query += " AND patronymic = $3"
	} else if patronymic == nil {
		query += " AND patronymic IS NULL"
	} else {
		query += " AND patronymic = ''"
	}

	var userID int64
	var err error

	if patronymic == nil {
		err = u.db.GetContext(ctx, &userID, query, name, surname)
	} else {
		err = u.db.GetContext(ctx, &userID, query, name, surname, patronymic)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, err
		}

		return -1, err
	}

	return userID, nil
}

func (u *UserRepo) CreateUserDataByName(ctx context.Context, individual models.Individual) (string, error) {
	individualData := individual.ToStorage()

	var lastInsertedID int64
	query := `
		INSERT INTO individuals (
			name, 
			surname, 
			patronymic, 
			age, 
			gender, 
			country_id
		) VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6
		) RETURNING id
	`

	err := u.db.QueryRowContext(
		ctx,
		query,
		individualData.Name,
		individualData.Surname,
		individualData.Patronymic,
		individualData.Age,
		individualData.Gender,
		individualData.CountryID,
	).Scan(&lastInsertedID)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(lastInsertedID)), nil
}

func (u *UserRepo) GetUserDataByID(ctx context.Context, idNumber int64) (models.UserData, error) {
	var userData models.UserData
	query := "SELECT name, surname, patronymic FROM individuals WHERE id = $1"

	err := u.db.GetContext(ctx, &userData, query, idNumber)
	if err != nil {
		return models.UserData{}, err
	}

	return userData, nil
}

func (u *UserRepo) DeleteAllUserDataByID(ctx context.Context, idNumber int64) error {
	query := "DELETE FROM individuals WHERE id = $1"

	command, err := u.db.ExecContext(ctx, query, idNumber)
	if err != nil {
		return err
	}

	row, err := command.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errlst.ErrUserNotFound
	}

	return nil
}

func (u *UserRepo) UpdateUserAllDataByID(ctx context.Context, individual models.Individual) error {
	individualData := individual.ToStorage()

	idNum, err := strconv.Atoi(individualData.ID)
	if err != nil {
		return err
	}

	query := `
		UPDATE individuals
		SET 
			name = $2,
			surname = $3,
			patronymic = $4,
			age = $5,
			gender = $6,
			country_id = $7,
			updated_at = $8
		WHERE 
			id = $1
		`

	_, err = u.db.ExecContext(
		ctx,
		query,
		idNum,
		individualData.Name,
		individualData.Surname,
		individualData.Patronymic,
		individualData.Age,
		individualData.Gender,
		individualData.CountryID,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) UpdateUserDataByID(ctx context.Context, idNumber int64, surname string, patronymic *string) error {
	query := "UPDATE individuals SET surname = $2, patronymic = $3, updated_at = $4  WHERE id = $1"

	_, err := u.db.ExecContext(ctx, query, idNumber, surname, patronymic, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetUserAllDataByID(ctx context.Context, idNumber int64) (*models.IndividualData, error) {
	query := "SELECT name, surname, patronymic, age, gender, country_id FROM individuals WHERE id = $1"
	var userData models.IndividualData

	err := u.db.GetContext(ctx, &userData, query, idNumber)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}
func (u *UserRepo) GetUserDataByFilter(
	ctx context.Context,
	filter *models.IndividualFilter,
	page *int,
	pageSize *int,
	sortBy *string,
) ([]*models.IndividualData, error) {
	query := "SELECT id, name, surname, patronymic, age, gender, country_id FROM individuals WHERE 1 = 1"
	var args []interface{} // Slice to hold the query arguments

	if filter != nil {
		if filter.Name != nil && *filter.Name != "" {
			query += " AND name = $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.Name)
		}

		if filter.Surname != nil && *filter.Surname != "" {
			query += " AND surname = $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.Surname)
		}

		if filter.Patronymic != nil && *filter.Patronymic != "" {
			query += " AND patronymic = $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.Patronymic)
		}

		if filter.AgeMin != nil && *filter.AgeMin > 0 {
			query += " AND age >= $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.AgeMin)
		}

		if filter.AgeMax != nil && *filter.AgeMax > 0 {
			query += " AND age <= $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.AgeMax)
		}

		if filter.Gender != nil && *filter.Gender != "" {
			query += " AND gender = $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.Gender)
		}

		if filter.CountryID != nil && *filter.CountryID != "" {
			query += " AND country_id = $" + strconv.Itoa(len(args)+1)
			args = append(args, *filter.CountryID)
		}
	}

	if sortBy != nil && *sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", *sortBy)
	}

	if page != nil && pageSize != nil {
		offset := (*page - 1) * *pageSize
		query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(args)+1, len(args)+2)
		args = append(args, offset, *pageSize)
	}

	// Execute the query
	rows, err := u.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var individuals []*models.IndividualData

	for rows.Next() {
		var individual models.IndividualData
		// Scan individual fields from the rows into the struct
		err := rows.Scan(
			&individual.ID,
			&individual.Name,
			&individual.Surname,
			&individual.Patronymic,
			&individual.Age,
			&individual.Gender,
			&individual.CountryID,
		)
		if err != nil {
			return nil, err
		}

		individuals = append(individuals, &individual)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return individuals, nil
}
