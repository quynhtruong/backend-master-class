package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername     = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidUserFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contains only letters, digits and underscore")
	}
	return nil
}

func ValidPassword(value string) error {
	return ValidateString(value, 3, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("input is not valid email")
	}
	return nil
}

func ValidateUserFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUserFullName(value) {
		return fmt.Errorf("must contains only letters or spaces")
	}
	return nil
}
