package options

import (
	"encoding/json"

	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
	"github.com/pachirode/iam_study/internal/pkg/server"
	"github.com/pachirode/iam_study/pkg/flags"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/utils/idutil"
)

type Options struct {
	GenericServerRunOptions *genericOptions.ServerRunOptions       `json:"server"   mapstructure:"server"`
	GRPCOptions             *genericOptions.GRPCOptions            `json:"grpc"     mapstructure:"grpc"`
	InsecureServing         *genericOptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericOptions.SecureServingOptions   `json:"secure"   mapstructure:"sercure"`
	MySQLOptions            *genericOptions.MySQLOptions           `json:"mysql"    mapstructure:"mysql"`
	RedisOptions            *genericOptions.RedisOptions           `json:"redis"    mapstructure:"redis"`
	JwtOptions              *genericOptions.JwtOptions             `json:"jwt"      mapstructure:"jwt"`
	Log                     *log.Options                           `json:"log"      mapstructure:"log"`
	FeatureOptions          *genericOptions.FeatureOptions         `json:"feature"  mapstructure:"feature"`
}

func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: genericOptions.NewServerRunOptions(),
		GRPCOptions:             genericOptions.NewGRPCOptions(),
		InsecureServing:         genericOptions.NewInsecureServingOptions(),
		SecureServing:           genericOptions.NewSecureServingOptions(),
		JwtOptions:              genericOptions.NewJwtOptions(),
		FeatureOptions:          genericOptions.NewFeatureOptions(),
		RedisOptions:            genericOptions.NewRedisOptions(),
		MySQLOptions:            genericOptions.NewMySQLOptions(),
		Log:                     log.NewOptions(),
	}
}

func (opt *Options) ApplyTo(config *server.Config) error {
	return nil
}

func (opt *Options) Flags() (nfs flags.NamedFlagSets) {
	opt.GenericServerRunOptions.AddFlags(nfs.GetFlagSet("generic"))
	opt.GRPCOptions.AddFlags(nfs.GetFlagSet("grpc"))
	opt.InsecureServing.AddFlags(nfs.GetFlagSet("insecure serving"))
	opt.SecureServing.AddFlags(nfs.GetFlagSet("sercure serving"))
	opt.FeatureOptions.AddFlags(nfs.GetFlagSet("features"))
	opt.JwtOptions.AddFlags(nfs.GetFlagSet("jwt"))
	opt.MySQLOptions.AddFlags(nfs.GetFlagSet("mysql"))
	opt.RedisOptions.AddFlags(nfs.GetFlagSet("redis"))
	opt.Log.AddFlags(nfs.GetFlagSet("logs"))

	return nfs
}

func (opt *Options) String() string {
	data, _ := json.Marshal(opt)

	return string(data)
}

func (opt *Options) Complete() error {
	if opt.JwtOptions.Key == "" {
		opt.JwtOptions.Key = idutil.NewSecretKey()
	}

	return opt.SecureServing.Complete()
}
