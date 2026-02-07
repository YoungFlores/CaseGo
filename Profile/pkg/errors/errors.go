package apperrors

import "errors"

var (
	ErrAlreadyExists = errors.New("resource_already_exists")
	ErrInternal      = errors.New("internal_server_error")
	ErrUsernameTaken = errors.New("username_is_already_taken")
	ErrEmailTaken    = errors.New("email_is_already_taken")
	ErrIsNotActive   = errors.New("user_is_not_active")
	ErrForbidden     = errors.New("havn't access")
)
