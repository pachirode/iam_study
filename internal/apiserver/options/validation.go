package options

func (opt *Options) Validate() []error {
	var errs []error

	errs = append(errs, opt.GenericServerRunOptions.Validate()...)
	errs = append(errs, opt.MySQLOptions.Validate()...)
	errs = append(errs, opt.Log.Validate()...)

	return errs
}
