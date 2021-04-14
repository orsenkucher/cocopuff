package log

import (
	"log"

	"go.uber.org/zap"
)

func New(service, deployment, version string, release bool) (*zap.SugaredLogger, error) {
	var (
		err    error
		logger *zap.Logger
	)

	if release {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("service", service),
		zap.String("dep", string(deployment)),
		zap.String("ver", version),
	)

	return logger.Sugar(), nil
}

func Abort(v ...interface{}) {
	log.Fatal(v...)
}

func Abortf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func Abortln(v ...interface{}) {
	log.Fatalln(v...)
}
