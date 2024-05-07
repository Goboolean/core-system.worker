package job

import "errors"

var TypeMismatchError = errors.New("type assertion: failed to convert type")
var NotFoundJob = errors.New("select job: not found desired spec")
