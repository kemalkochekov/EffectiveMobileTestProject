package models

type AgifyResponse struct {
	Count int    `json:"count"`
	Age   int    `json:"age"`
	Name  string `json:"name"`
}
type NationalizeResponse struct {
	Count   int                `json:"count"`
	Name    string             `json:"name"`
	Country []*NationalityInfo `json:"country"`
}
type NationalityInfo struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type GenderizeResponse struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type UserData struct {
	ID         int64   `db:"id"`
	Name       string  `db:"name"`
	Surname    string  `db:"surname"`
	Patronymic *string `db:"patronymic"`
}
type IndividualFilterData struct {
	Name       *string `db:"name"`
	Surname    *string `db:"surname"`
	Patronymic *string `db:"patronymic"`
	AgeMin     *int    `db:"ageMin"`
	AgeMax     *int    `db:"ageMax"`
	Gender     *string `db:"gender"`
	CountryID  *string `db:"country_id"`
}
type IndividualData struct {
	ID         string  `db:"id"`
	Name       string  `db:"name"`
	Surname    string  `db:"surname"`
	Patronymic *string `db:"patronymic"`
	Age        *int    `db:"age"`
	Gender     *string `db:"gender"`
	CountryID  *string `db:"country_id"`
}

func (i *Individual) ToStorage() IndividualData {
	return IndividualData{
		ID:         i.ID,
		Name:       i.Name,
		Surname:    i.Surname,
		Patronymic: i.Patronymic,
		Age:        i.Age,
		Gender:     i.Gender,
		CountryID:  i.CountryID,
	}
}
func (i *IndividualData) ToServer() *Individual {
	return &Individual{
		ID:         i.ID,
		Name:       i.Name,
		Surname:    i.Surname,
		Patronymic: i.Patronymic,
		Age:        i.Age,
		Gender:     i.Gender,
		CountryID:  i.CountryID,
	}
}

func (i *IndividualFilterData) ToStorage() IndividualFilterData {
	return IndividualFilterData{
		Name:       i.Name,
		Surname:    i.Surname,
		Patronymic: i.Patronymic,
		AgeMin:     i.AgeMin,
		AgeMax:     i.AgeMax,
		Gender:     i.Gender,
		CountryID:  i.CountryID,
	}
}
