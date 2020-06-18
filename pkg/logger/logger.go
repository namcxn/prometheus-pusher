package logger

import (
	"github.com/op/go-logging"
	"go.uber.org/zap"
	"os"
)


func getZapLogger(module string) *zap.SugaredLogger {
	log, _ := zap.NewProduction()
	if module == "" {
		return log.Sugar()
	}
	return log.Sugar().With("module", module)
}

var log *logging.Logger

func GetLogger(module string) *logging.Logger{
	log = logging.MustGetLogger(module)
	var format = logging.MustStringFormatter(
		`%{color} - %{time:15:04:05.000} - %{shortfunc} - %{shortfile} - %{level:.4s} %{color:reset}:   %{message}`,
	)
	logBackend := logging.NewLogBackend(os.Stdout, "*", 0)
	logBackendFormatter := logging.NewBackendFormatter(logBackend, format)
	logging.SetBackend(logBackendFormatter)
	logging.SetLevel(logging.INFO, "")
	return log
}