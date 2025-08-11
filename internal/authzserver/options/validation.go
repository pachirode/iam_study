package options

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.GenericServerRunOptions.Validate()...)
	errs = append(errs, o.InsecureServing.Validate()...)
	errs = append(errs, o.SecureServing.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.AnalyticsOptions.Validate()...)

	return errs
}
