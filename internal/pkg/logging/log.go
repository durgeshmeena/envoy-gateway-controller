package logging

import (
	"flag"

	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func InitLogger() logr.Logger {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts))
	logf.SetLogger(logger)
	return logger
}

// func NewLoggerWithName(name string) logr.Logger {
// 	return logf.Log.WithName(name)
// }