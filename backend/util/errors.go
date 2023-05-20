package util

import "github.com/secureworks/errors"

const (
	UserFacingTagKey = "userFacing"
)

var (
	UserFacingTag       = errors.Tag("userFacing")
	ErrBadRequest       = errors.New("Bad request", errors.Meta("code", 400))
	ErrPermissionDenied = errors.New("Permission denied", errors.Meta("code", 401))
)
