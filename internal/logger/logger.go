package logger

import (
	"go.uber.org/zap"

	"github.com/Nexadis/GophKeeper/internal/config"
)

var logger *zap.SugaredLogger

func Init(c *config.LogConfig) error {
	lvl, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = c.Outputs
	cfg.Level = lvl
	if c.Encoding != "" {
		cfg.Encoding = c.Encoding

	}
	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	logger = l.Sugar()
	return nil
}

func Info(args ...interface{}) {
	logger.Info(args...)
}
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}
