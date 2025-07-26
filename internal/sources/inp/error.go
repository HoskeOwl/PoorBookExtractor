package inp

import "errors"

var ErrDeleted = errors.New("deleted book")
var ErrInvalidSize = errors.New("invalid size")
var ErrInvalidDate = errors.New("invalid date")
