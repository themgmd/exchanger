package zaplog

import (
	"go.uber.org/zap"
	"onemgvv/exchanger/pkg/logger"
)

var (
	appPath  = "/var/log/healfy_backend/app.log"
	httpPath = "/var/log/healfy_backend/http.log"

	AppLogger  = New(appPath)
	HttpLogger = New(httpPath)
)

func New(path string) *zap.SugaredLogger {
	return logger.New(path)
}
