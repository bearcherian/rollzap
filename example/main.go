package main

import (
	"github.com/bearcherian/rollzap"
	rollbar "github.com/rollbar/rollbar-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Initialize rollbar with your token and optional environment flag
	rollbar.SetToken("MY_ROLL_BAR_TOKEN")
	rollbar.SetEnvironment("production")

	// create a new core that sends zapcore.WarnLevel and above messages to Rollbar
	rollbarCore := rollzap.NewRollbarCore(zapcore.WarnLevel)

	// create a new Zap logger
	logger, _ := zap.NewProduction()

	// Wrap a NewTee to send log messages to both your main logger and to rollbar
	logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, rollbarCore)
	}))

	// This message will only go to the main logger
	logger.Info("Rollbar Core teed up", zap.String("foo", "bar"))

	// This warning will go to both the main logger and to Rollbar. the 'foo' field will appear in rollbar as 'custom.foo'
	logger.Warn("Warning message with fields", zap.String("foo", "bar"))
}
