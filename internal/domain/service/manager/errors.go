package manager

import "errors"



var ErrRegisterFailed = errors.New("failed to register worker")

var ErrInvalidTaskType = errors.New("invalid task type")