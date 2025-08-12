package options

import (
	"encoding/json"

	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	"github.com/pachirode/iam_study/internal/pump/analytics"
	"github.com/pachirode/iam_study/pkg/flags"
	"github.com/pachirode/iam_study/pkg/log"
)

type PumpConfig struct {
	Type                  string                     `json:"type"                    mapstructure:"type"`
	Filters               analytics.AnalyticsFilters `json:"filters"                 mapstructure:"filters"`
	Timeout               int                        `json:"timeout"                 mapstructure:"timeout"`
	OmitDetailedRecording bool                       `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	Meta                  map[string]interface{}     `json:"meta"                    mapstructure:"meta"`
}

type Options struct {
	PurgeDelay            int                          `json:"purge-delay"             mapstructure:"purge-delay"`
	Pumps                 map[string]PumpConfig        `json:"pumps"                   mapstructure:"pumps"`
	HealthCheckPath       string                       `json:"health-check-path"       mapstructure:"health-check-path"`
	HealthCheckAddress    string                       `json:"health-check-address"    mapstructure:"health-check-address"`
	OmitDetailedRecording bool                         `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	RedisOptions          *genericOptions.RedisOptions `json:"redis"                   mapstructure:"redis"`
	Log                   *log.Options                 `json:"log"                     mapstructure:"log"`
}

func NewOptions() *Options {
	s := Options{
		PurgeDelay: 10,
		Pumps: map[string]PumpConfig{
			"csv": {
				Type: "csv",
				Meta: map[string]interface{}{
					"csv_dir": "./analytics-data",
				},
			},
		},
		HealthCheckPath:    "healthz",
		HealthCheckAddress: "0.0.0.0:7070",
		RedisOptions:       genericOptions.NewRedisOptions(),
		Log:                log.NewOptions(),
	}

	return &s
}

func (o *Options) Flags() (nfs flags.NamedFlagSets) {
	o.RedisOptions.AddFlags(nfs.GetFlagSet("redis"))
	o.Log.AddFlags(nfs.GetFlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := nfs.GetFlagSet("misc")
	fs.IntVar(&o.PurgeDelay, "purge-delay", o.PurgeDelay, ""+
		"This setting the purge delay (in seconds) when purge the data from Redis to MongoDB or other data stores.")
	fs.StringVar(&o.HealthCheckPath, "health-check-path", o.HealthCheckPath, ""+
		"Specifies liveness health check request path.")
	fs.StringVar(&o.HealthCheckAddress, "health-check-address", o.HealthCheckAddress, ""+
		"Specifies liveness health check bind address.")
	fs.BoolVar(&o.OmitDetailedRecording, "omit-detailed-recording", o.OmitDetailedRecording, ""+
		"Setting this to true will avoid writing policy fields for each authorization request in pumps.")

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
