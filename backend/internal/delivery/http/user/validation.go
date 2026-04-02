package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errFIORequired            = errors.New("fio is required")
	errTelephoneNumberNeeded  = errors.New("telephone number is required")
	errInvalidTelephoneNumber = errors.New("invalid telephone number")
	phonePattern              = regexp.MustCompile(`^(8[0-9]{10}|\+7[0-9]{10})$`)
	errFIOTooLong             = errors.New("fio is too long")
	errTelephoneNumberTooLong = errors.New("telephone number is too long")
	errCityTooLong            = errors.New("city is too long")
	errUserLoginTooLong       = errors.New("user login is too long")
	errUserPasswordTooLong    = errors.New("user password is too long")
	errStatusTooLong          = errors.New("status is too long")
	errUserDescriptionTooLong = errors.New("user description is too long")
)

const (
	maxFIOLength             = 255
	maxTelephoneNumberLength = 20
	maxCityLength            = 50
	maxUserLoginLength       = 50
	maxUserPasswordLength    = 255
	maxStatusLength          = 20
	maxUserDescriptionLength = 1000
)

func validateTelephoneNumber(phone string) error {
	normalized := strings.TrimSpace(phone)
	if !phonePattern.MatchString(normalized) {
		return errInvalidTelephoneNumber
	}
	return nil
}

func validateUserRequiredFields(fio string, telephoneNumber string) error {
	switch {
	case strings.TrimSpace(fio) == "":
		return errFIORequired
	case strings.TrimSpace(telephoneNumber) == "":
		return errTelephoneNumberNeeded
	default:
		return nil
	}
}

func validateUserDescription(description string) error {
	if len(description) > maxUserDescriptionLength {
		return errUserDescriptionTooLong
	}
	return nil
}

func validateUserFieldsMaxLength(
	fio string,
	telephoneNumber string,
	city string,
	userLogin string,
	userPassword string,
	status string,
	userDescription string,
) error {
	switch {
	case len(fio) > maxFIOLength:
		return errFIOTooLong
	case len(telephoneNumber) > maxTelephoneNumberLength:
		return errTelephoneNumberTooLong
	case len(city) > maxCityLength:
		return errCityTooLong
	case len(userLogin) > maxUserLoginLength:
		return errUserLoginTooLong
	case len(userPassword) > maxUserPasswordLength:
		return errUserPasswordTooLong
	case len(status) > maxStatusLength:
		return errStatusTooLong
	case validateUserDescription(userDescription) != nil:
		return errUserDescriptionTooLong
	default:
		return nil
	}
}

func userValidationErrorMessage(err error) string {
	switch {
	case errors.Is(err, errFIORequired):
		return "Поле ФИО обязательно"
	case errors.Is(err, errTelephoneNumberNeeded):
		return "Поле номера телефона обязательно"
	case errors.Is(err, errInvalidTelephoneNumber):
		return "Неверный формат номера телефона"
	case errors.Is(err, errFIOTooLong):
		return "Слишком длинное ФИО"
	case errors.Is(err, errTelephoneNumberTooLong):
		return "Слишком длинный номер телефона"
	case errors.Is(err, errCityTooLong):
		return "Слишком длинное название города"
	case errors.Is(err, errUserLoginTooLong):
		return "Слишком длинный логин"
	case errors.Is(err, errUserPasswordTooLong):
		return "Слишком длинный пароль"
	case errors.Is(err, errStatusTooLong):
		return "Слишком длинный статус"
	case errors.Is(err, errUserDescriptionTooLong):
		return "Слишком длинное описание пользователя"
	default:
		return "Некорректные данные пользователя"
	}
}
