package domain

import "errors"

var (
	// ErrInternalServerError will throw if any internal server error occurs
	ErrInternalServerError = errors.New("internal Server Error")

	// ErrNotFound will throw if the requested item does not exist
	ErrNotFound = errors.New("your requested item is not found")

	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("item already exists")

	// ErrBadParamInput will throw if the given request body or params are not valid
	ErrBadParamInput = errors.New("given parameter is not valid")

	// ErrUnauthorized will throw if the user is not authorized
	ErrUnauthorized = errors.New("user is not authorized")

	// ErrInvalidCredentials will throw if login or password is incorrect
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrAccountLocked will throw if the user’s account is locked
	ErrAccountLocked = errors.New("account is locked")

	// ErrSessionExpired will throw if the session token is expired
	ErrSessionExpired = errors.New("session has expired")

	// ErrInsufficientFunds will throw if the user does not have enough funds
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrInvalidToken will throw if the provided token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrForbidden will throw if the user does not have permission for the action
	ErrForbidden = errors.New("access is forbidden")

	// ErrTooManyRequests will throw if the rate limit is exceeded
	ErrTooManyRequests = errors.New("too many requests, please try again later")
)
