package analytics

import (
	"github.com/spf13/pflag"
)

type AnalyticsOptions struct {
	RecordsBufferSize       uint64 `json:"records-buffer-size" mapstructure:"records-buffer-size"`
	Enable                  bool   `json:"enable" mapstructure:"enable"`
	EnableDetailedRecording bool   `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"`
}

func NewAnalyticsOptions() *AnalyticsOptions {
	return &AnalyticsOptions{
		Enable:                  true,
		RecordsBufferSize:       2000,
		EnableDetailedRecording: true,
	}
}

func (o *AnalyticsOptions) Validate() []error {
	return []error{}
}

func (o *AnalyticsOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.Enable, "analytics.enable", o.Enable, "Set analytics is enable")
	fs.Uint64Var(&o.RecordsBufferSize, "analytics.records-buffer-size", o.RecordsBufferSize, "Set analytics records buffer size")
	fs.BoolVar(&o.EnableDetailedRecording, "analytics.enable-detailed-recording", o.EnableDetailedRecording, "Set enable detailed recording")
}
