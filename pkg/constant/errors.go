package constant

import "errors"

// internal
var (
	ErrCaptchaCode          = errors.New("ErrCaptchaCode")
	ErrUserDisabled         = errors.New("ErrUserDisabled")
	ErrWhiteIpList          = errors.New("ErrWhiteIpList")
	ErrInvalidParams        = errors.New("ErrInvalidParams")
	ErrFailedAuthentication = errors.New("ErrFailedAuthentication")
	ErrUnauthorized         = errors.New("ErrUnauthorized")
	ErrErrForbidden         = errors.New("ErrForbidden")
	ErrCreated              = errors.New("ErrCreated")
	ErrUpdated              = errors.New("ErrUpdated")
	ErrDeleted              = errors.New("ErrDeleted")
	ErrRecordNotFound       = errors.New("ErrRecordNotFound")
	ErrUpdateForbidden      = errors.New("ErrUpdateForbidden")
	ErrInternalServer       = errors.New("InternalServer")
)
