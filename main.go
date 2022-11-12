package main

import (
	"exampleproject/expr"
	"exampleproject/log"
	"time"

	"go.uber.org/zap"
)

func main() {

	log.Log("hello", log.CategoryID(123), log.Exprs(expr.CategoryIDEquals(234)))
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "URL",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", "url")
}
