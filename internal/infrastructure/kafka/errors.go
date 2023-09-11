package kafka

import "errors"


var ErrDeadlineSettingRequired = errors.New("deadline setting is required on context")

var ErrFailedToFlush = errors.New("failed to flush")

var ErrReceivedMsgIsNotProtoMessage = errors.New("received message is not proto message")

var ErrTopicAlreadySubscribed = errors.New("topic is already subscribed")