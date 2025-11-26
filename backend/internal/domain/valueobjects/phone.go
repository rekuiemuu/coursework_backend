package valueobjects

import (
	"fmt"
	"regexp"
)

type Phone struct {
	Value string
}

func NewPhone(value string) (*Phone, error) {
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	cleaned := regexp.MustCompile(`[\s\-\(\)]`).ReplaceAllString(value, "")

	if !phoneRegex.MatchString(cleaned) {
		return nil, fmt.Errorf("invalid phone format")
	}
	return &Phone{Value: cleaned}, nil
}

func (p *Phone) String() string {
	return p.Value
}
