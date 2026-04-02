package user

import (
	"strings"
	"testing"
)

func TestValidateTelephoneNumber(t *testing.T) {
	t.Parallel()

	validNumbers := []string{
		"+79998887766",
		"89998887766",
	}

	for _, number := range validNumbers {
		number := number
		t.Run("valid_"+number, func(t *testing.T) {
			t.Parallel()

			if err := validateTelephoneNumber(number); err != nil {
				t.Fatalf("expected valid number %q, got error: %v", number, err)
			}
		})
	}

	invalidNumbers := []string{
		"",
		"12345",
		"phone",
		"+7999 8887766",
		"+7(999)8887766",
		"+7-999-888-77-66",
		"79998887766",
		"+89998887766",
		"+7-999-888-77-66-00-11-22",
		"+7(999)ABC-77-66",
	}

	for _, number := range invalidNumbers {
		number := number
		t.Run("invalid_"+number, func(t *testing.T) {
			t.Parallel()

			if err := validateTelephoneNumber(number); err == nil {
				t.Fatalf("expected invalid number %q to fail validation", number)
			}
		})
	}
}

func TestValidateUserDescription(t *testing.T) {
	t.Parallel()

	validDescription := strings.Repeat("a", maxUserDescriptionLength)
	if err := validateUserDescription(validDescription); err != nil {
		t.Fatalf("expected description with %d chars to be valid, got: %v", maxUserDescriptionLength, err)
	}

	invalidDescription := strings.Repeat("a", maxUserDescriptionLength+1)
	if err := validateUserDescription(invalidDescription); err == nil {
		t.Fatalf("expected description with %d chars to be invalid", maxUserDescriptionLength+1)
	}
}

func TestValidateUserFieldsMaxLength(t *testing.T) {
	t.Parallel()

	validErr := validateUserFieldsMaxLength(
		strings.Repeat("a", maxFIOLength),
		strings.Repeat("1", maxTelephoneNumberLength),
		strings.Repeat("b", maxCityLength),
		strings.Repeat("c", maxUserLoginLength),
		strings.Repeat("d", maxUserPasswordLength),
		strings.Repeat("e", maxStatusLength),
		strings.Repeat("f", maxUserDescriptionLength),
	)
	if validErr != nil {
		t.Fatalf("expected max-length values to be valid, got: %v", validErr)
	}

	tests := []struct {
		name string
		err  error
		in   [7]string
	}{
		{
			name: "fio too long",
			err:  errFIOTooLong,
			in: [7]string{
				strings.Repeat("a", maxFIOLength+1),
				"89998887766",
				"City",
				"login",
				"password",
				"active",
				"desc",
			},
		},
		{
			name: "telephone too long",
			err:  errTelephoneNumberTooLong,
			in: [7]string{
				"fio",
				strings.Repeat("1", maxTelephoneNumberLength+1),
				"City",
				"login",
				"password",
				"active",
				"desc",
			},
		},
		{
			name: "city too long",
			err:  errCityTooLong,
			in: [7]string{
				"fio",
				"89998887766",
				strings.Repeat("b", maxCityLength+1),
				"login",
				"password",
				"active",
				"desc",
			},
		},
		{
			name: "login too long",
			err:  errUserLoginTooLong,
			in: [7]string{
				"fio",
				"89998887766",
				"City",
				strings.Repeat("c", maxUserLoginLength+1),
				"password",
				"active",
				"desc",
			},
		},
		{
			name: "password too long",
			err:  errUserPasswordTooLong,
			in: [7]string{
				"fio",
				"89998887766",
				"City",
				"login",
				strings.Repeat("d", maxUserPasswordLength+1),
				"active",
				"desc",
			},
		},
		{
			name: "status too long",
			err:  errStatusTooLong,
			in: [7]string{
				"fio",
				"89998887766",
				"City",
				"login",
				"password",
				strings.Repeat("e", maxStatusLength+1),
				"desc",
			},
		},
		{
			name: "description too long",
			err:  errUserDescriptionTooLong,
			in: [7]string{
				"fio",
				"89998887766",
				"City",
				"login",
				"password",
				"active",
				strings.Repeat("f", maxUserDescriptionLength+1),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := validateUserFieldsMaxLength(
				tt.in[0],
				tt.in[1],
				tt.in[2],
				tt.in[3],
				tt.in[4],
				tt.in[5],
				tt.in[6],
			)
			if got != tt.err {
				t.Fatalf("expected %v, got %v", tt.err, got)
			}
		})
	}
}

func TestValidateUserRequiredFields(t *testing.T) {
	t.Parallel()

	if err := validateUserRequiredFields("Ivan Ivanov", "89998887766"); err != nil {
		t.Fatalf("expected required fields to pass validation, got: %v", err)
	}

	if err := validateUserRequiredFields("", "89998887766"); err != errFIORequired {
		t.Fatalf("expected errFIORequired, got: %v", err)
	}

	if err := validateUserRequiredFields("Ivan Ivanov", "   "); err != errTelephoneNumberNeeded {
		t.Fatalf("expected errTelephoneNumberNeeded, got: %v", err)
	}
}
