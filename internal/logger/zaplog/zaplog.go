package zaplog

import (
	"exchanger/pkg/logger"
	"go.uber.org/zap"
)

var (
	//appPath  = "/var/log/healfy_backend/app.log"
	//httpPath = "/var/log/healfy_backend/http.log"

	appPath  = "logs/app.log"
	httpPath = "logs/http.log"

	AppLogger  = New(appPath)
	HttpLogger = New(httpPath)
)

func New(path string) *zap.SugaredLogger {
	return logger.New(path)
}
