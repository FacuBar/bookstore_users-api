package ports

import "go.uber.org/zap"

type UserLogger interface {
	Error(msg string, err error, tags ...zap.Field)
	Info(msg string, tags ...zap.Field)
}
