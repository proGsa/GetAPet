package user

import "testing"

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
