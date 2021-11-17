package logger

import (
	"sync"

	"github.com/FacuBar/bookstore_users-api/pkg/core/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type userLogger struct {
	log *zap.Logger
}

var (
	onceUsersLogger sync.Once
	instanceLogger  *userLogger
)

func NewUserLogger() ports.UserLogger {
	onceUsersLogger.Do(func() {
		logConfig := zap.Config{
			OutputPaths: []string{"stdout"},
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				LevelKey:     "level",
				TimeKey:      "time",
				MessageKey:   "msg",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeLevel:  zapcore.LowercaseLevelEncoder,
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}

		log, err := logConfig.Build()
		if err != nil {
			panic(err)
		}

		instanceLogger = &userLogger{
			log: log,
		}
	})
	return instanceLogger
}

func (l *userLogger) Info(msg string, tags ...zap.Field) {
	l.log.Info(msg, tags...)
	l.log.Sync()
}

func (l *userLogger) Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Error(msg, tags...)
	l.log.Sync()
}
