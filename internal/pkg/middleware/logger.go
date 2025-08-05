package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-isatty"

	"github.com/pachirode/iam_study/pkg/log"
)

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency - param.Latency%time.Second
	}

	return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
		param.TimeStamp.Format("2001-01-01 - 10:00:00"),
		statusColor, param.StatusCode, resetColor,
		param.ClientIP,
		param.Latency,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func GetDefaultLogFormatter() gin.LogFormatter {
	return func(params gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if params.IsOutputColor() {
			statusColor = params.StatusCodeColor()
			methodColor = params.MethodColor()
			resetColor = params.ResetColor()
		}

		if params.Latency > time.Minute {
			params.Latency -= params.Latency % time.Second
		}

		return fmt.Sprintf("%s%3d%s - [%s] \"%v %s%s%s %s\" %s",
			statusColor, params.StatusCode, resetColor,
			params.ClientIP,
			params.Latency,
			methodColor, params.Method, resetColor,
			params.Path,
			params.ErrorMessage,
		)
	}
}

func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skipPaths []string) gin.LoggerConfig {
	if formatter == nil {
		formatter = GetDefaultLogFormatter()
	}

	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

func LoggerWithConfig(config gin.LoggerConfig) gin.HandlerFunc {
	formatter := config.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := config.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notLogged := config.SkipPaths
	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	if isTerm {
		gin.ForceConsoleColor()
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip := make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		if _, ok := skip[path]; !ok {
			params := gin.LogFormatterParams{
				Request: ctx.Request,
				Keys:    ctx.Keys,
			}

			params.TimeStamp = time.Now()
			params.Latency = params.TimeStamp.Sub(start)

			params.ClientIP = ctx.ClientIP()
			params.Method = ctx.Request.Method
			params.StatusCode = ctx.Writer.Status()
			params.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()

			params.BodySize = ctx.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			params.Path = path

			log.L(ctx).Info(formatter(params))
		}
	}
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(GetLoggerConfig(nil, nil, nil))
}
