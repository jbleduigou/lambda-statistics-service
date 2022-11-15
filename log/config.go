package log

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(ctx context.Context) {
	// Retrieve AWS Request ID
	lc, _ := lambdacontext.FromContext(ctx)
	requestID := lc.AwsRequestID
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(getLogLevel()),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"request-id": requestID},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func getLogLevel() zapcore.Level {
	lvl, found := os.LookupEnv("LOG_LEVEL")
	if found {
		var l zapcore.Level
		if err := l.Set(lvl); err == nil {
			return l
		}
	}
	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
		return zap.DebugLevel
	}
	if env == "staging" {
		return zap.InfoLevel
	}
	return zap.WarnLevel
}
