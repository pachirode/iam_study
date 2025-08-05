package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	flagDevelopment       = "log.development"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagEnableColor       = "log.enable-color"
	flagFormat            = "log.format"
	flagLevel             = "log.level"
	flagName              = "log.name"
	flagOutputPaths       = "log.output-paths"

	consoleFormat = "console"
	jsonFormat    = "json"
)

type Options struct {
	Development       bool     `json:"development"        mapstructure:"development"`
	DisableCaller     bool     `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	EnableColor       bool     `json:"enable-color"       mapstructure:"enable-color"`
	Format            string   `json:"format"             mapstructure:"format"`
	Level             string   `json:"level"              mapstructure:"level"`
	Name              string   `json:"name"               mapstructure:"name"`
	OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
}

func NewOptions() *Options {
	return &Options{
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		ErrorOutputPaths:  []string{"stderr"},
		EnableColor:       false,
		Format:            consoleFormat,
		Level:             zapcore.InfoLevel.String(),
		OutputPaths:       []string{"stdout"},
	}
}

func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format %q", o.Format))
	}

	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log Level.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(
		&o.DisableStacktrace,
		flagDisableStacktrace,
		o.DisableStacktrace,
		"Disable the log to record a stack trace for message or above panic level.",
	)
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output format, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi color in plain logs")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log")
	fs.BoolVar(&o.Development, flagDevelopment, o.Development, "Enable development mode")
	fs.StringVar(&o.Name, flagName, o.Name, "The name of logger")
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Options) Build() error {
	var zapLevel zapcore.Level

	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
	}

	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}

	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}
