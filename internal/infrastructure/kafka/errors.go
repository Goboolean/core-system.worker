package kafka

import "errors"


var ErrDeadlineSettingRequired = errors.New("deadline setting is required on context")

var ErrFailedToFlush = errors.New("failed to flush")
