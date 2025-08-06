package options

func (opt *Options) Validate() []error {
	var errs []error

	errs = append(errs, opt.GenericServerRunOptions.Validate()...)
	errs = append(errs, opt.GRPCOptions.Validate()...)
	errs = append(errs, opt.InsecureServing.Validate()...)
	errs = append(errs, opt.SecureServing.Validate()...)
	errs = append(errs, opt.JwtOptions.Validate()...)
	errs = append(errs, opt.FeatureOptions.Validate()...)
	errs = append(errs, opt.MySQLOptions.Validate()...)
	errs = append(errs, opt.Log.Validate()...)

	return errs
}
