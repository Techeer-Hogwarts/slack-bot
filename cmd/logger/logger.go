package logger

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func JsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			if params.ErrorMessage != "" {
				log["error"] = params.ErrorMessage
			}

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
