package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrValidation     = errors.New("validation error")
	ErrDuplicate      = errors.New("duplicate entry")
	ErrInvalidState   = errors.New("invalid state transition")
	ErrQuestFull      = errors.New("quest has reached maximum tasks")
	ErrAlreadyStarted = errors.New("quest already started by this user")
)

type NotFoundError struct {
	Entity string
	ID     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Entity, e.ID)
}

func (e *NotFoundError) Unwrap() error {
	return ErrNotFound
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation of field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

type DuplicateError struct {
	Entity string
	Field  string
	Value  string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("%s from %s='%s' already exists", e.Entity, e.Field, e.Value)
}

func (e *DuplicateError) Unwrap() error {
	return ErrDuplicate
}
