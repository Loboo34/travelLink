// logger.go
package utils

import (
    "log"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(production bool) {
    var err error

    if production {
        Logger, err = zap.NewProduction()
    } else {
        // development logger — human readable, colored output
        config := zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
        Logger, err = config.Build()
    }

    if err != nil {
        log.Fatalf("failed to initialize logger: %v", err)
    }
}