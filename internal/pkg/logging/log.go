package logging

import (
	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var Log logr.Logger

func InitLogger(opts *zap.Options) {
	Log = zap.New(zap.UseFlagOptions(opts))
	logf.SetLogger(Log)
}