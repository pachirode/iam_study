package options

import (
	"encoding/json"

	"github.com/pachirode/iam_study/internal/authzserver/analytics"
	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	"github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/flags"
	"github.com/pachirode/iam_study/pkg/log"
)

type Options struct {
	RPCServer               string                                 `json:"rpcserver"      mapstructure:"rpcserver"`
	ClientCA                string                                 `json:"client-ca-file" mapstructure:"client-ca-file"`
	GenericServerRunOptions *genericOptions.ServerRunOptions       `json:"server"         mapstructure:"server"`
	InsecureServing         *genericOptions.InsecureServingOptions `json:"insecure"       mapstructure:"insecure"`
	SecureServing           *genericOptions.SecureServingOptions   `json:"secure"         mapstructure:"secure"`
	RedisOptions            *genericOptions.RedisOptions           `json:"redis"          mapstructure:"redis"`
	FeatureOptions          *genericOptions.FeatureOptions         `json:"feature"        mapstructure:"feature"`
	Log                     *log.Options                           `json:"log"            mapstructure:"log"`
	AnalyticsOptions        *analytics.AnalyticsOptions            `json:"analytics"      mapstructure:"analytics"`
}

func NewOptions() *Options {
	o := Options{
		RPCServer:               "127.0.0.1:8081",
		ClientCA:                "",
		GenericServerRunOptions: genericOptions.NewServerRunOptions(),
		InsecureServing:         genericOptions.NewInsecureServingOptions(),
		SecureServing:           genericOptions.NewSecureServingOptions(),
		RedisOptions:            genericOptions.NewRedisOptions(),
		FeatureOptions:          genericOptions.NewFeatureOptions(),
		Log:                     log.NewOptions(),
		AnalyticsOptions:        analytics.NewAnalyticsOptions(),
	}

	return &o
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (nfs flags.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(nfs.GetFlagSet("generic"))
	o.AnalyticsOptions.AddFlags(nfs.GetFlagSet("analytics"))
	o.RedisOptions.AddFlags(nfs.GetFlagSet("redis"))
	o.FeatureOptions.AddFlags(nfs.GetFlagSet("features"))
	o.InsecureServing.AddFlags(nfs.GetFlagSet("insecure serving"))
	o.SecureServing.AddFlags(nfs.GetFlagSet("secure serving"))
	o.Log.AddFlags(nfs.GetFlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := nfs.GetFlagSet("misc")
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "The address of iam rpc server. "+
		"The rpc server can provide all the secrets and policies to use.")
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")

	return nfs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	return o.SecureServing.Complete()
}
