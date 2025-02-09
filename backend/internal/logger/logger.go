package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
	}
}
