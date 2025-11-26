package valueobjects

import (
	"fmt"
	"regexp"
)

type Email struct {
	Value string
}

func NewEmail(value string) (*Email, error) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return nil, fmt.Errorf("invalid email format")
	}
	return &Email{Value: value}, nil
}

func (e *Email) String() string {
	return e.Value
}
