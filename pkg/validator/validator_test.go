package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestInput struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

func TestValidate_Success(t *testing.T) {
	Init()

	input := TestInput{
		Name:  "John",
		Email: "john@example.com",
	}

	errors := Validate(&input)

	assert.Empty(t, errors)
}

func TestValidate_RequiredField(t *testing.T) {
	Init()

	input := TestInput{
		Name:  "",
		Email: "john@example.com",
	}

	errors := Validate(&input)

	assert.Len(t, errors, 1)
	assert.Equal(t, "name", errors[0].Field)
	assert.Equal(t, "required", errors[0].Tag)
}

func TestValidate_InvalidEmail(t *testing.T) {
	Init()

	input := TestInput{
		Name:  "John",
		Email: "invalid-email",
	}

	errors := Validate(&input)

	assert.Len(t, errors, 1)
	assert.Equal(t, "email", errors[0].Field)
	assert.Equal(t, "email", errors[0].Tag)
}

func TestValidate_MinLength(t *testing.T) {
	Init()

	input := TestInput{
		Name:  "J",
		Email: "john@example.com",
	}

	errors := Validate(&input)

	assert.Len(t, errors, 1)
	assert.Equal(t, "name", errors[0].Field)
	assert.Equal(t, "min", errors[0].Tag)
}

func TestValidate_MultipleErrors(t *testing.T) {
	Init()

	input := TestInput{
		Name:  "",
		Email: "invalid",
	}

	errors := Validate(&input)

	assert.Len(t, errors, 2)
}