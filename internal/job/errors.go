package job

import "errors"

var ErrTypeMismatch = errors.New("type assertion: failed to convert type")
var ErrNotFoundJob = errors.New("select job: not found desired spec")
