package errlst

import "errors"

var (
	ErrEmptyNationalize  = errors.New("unable to get nationalize information")
	ErrEmptyAgify        = errors.New("unable to get agify information")
	ErrEmptyGenderize    = errors.New("unable to get genderize information")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
