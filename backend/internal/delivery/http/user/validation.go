package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errInvalidTelephoneNumber = errors.New("invalid telephone number")
	phonePattern              = regexp.MustCompile(`^(8[0-9]{10}|\+7[0-9]{10})$`)
)

func validateTelephoneNumber(phone string) error {
	normalized := strings.TrimSpace(phone)

	if !phonePattern.MatchString(normalized) {
		return errInvalidTelephoneNumber
	}

	return nil
}
