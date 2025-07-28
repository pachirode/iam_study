package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/flags"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/term"
	"github.com/pachirode/iam_study/pkg/version"
	"github.com/pachirode/iam_study/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	progressMessage = color.GreenString("==>")
	usageTemplate   = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

type App struct {
	basename    string
	name        string
	description string
	options     ClipOptions
	runFunc     RunFunc
	silence     bool
	noVersion   bool
	noConfig    bool
	commands    []*Command
	args        cobra.PositionalArgs
	cmd         *cobra.Command
}

func (app *App) AddCommand(cmd *Command) {
	app.commands = append(app.commands, cmd)
}

func (app *App) AddCommands(cmds ...*Command) {
	app.commands = append(app.commands, cmds...)
}

type Option func(*App)

func WithOptions(opts ClipOptions) Option {
	return func(app *App) {
		app.options = opts
	}
}

type RunFunc func(basename string) error

func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runFunc = run
	}
}

func WithDescription(desc string) Option {
	return func(app *App) {
		app.description = desc
	}
}

func WithSilence() Option {
	return func(app *App) {
		app.silence = true
	}
}

func WithNoVersion() Option {
	return func(app *App) {
		app.noVersion = true
	}
}

func WithNoConfig() Option {
	return func(app *App) {
		app.noConfig = true
	}
}

func WithValidArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

func WithDefaultValidArgs() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func NewApp(name string, basename string, opts ...Option) *App {
	app := &App{
		name:     name,
		basename: basename,
	}

	for _, opt := range opts {
		opt(app)
	}

	app.buildCommand()

	return app
}

func (app *App) buildCommand() {
	flags.InitFlags()

	cobraComand := cobra.Command{
		Use:   FormatBasename(app.basename),
		Short: app.name,
		Long:  app.description,
		// stop printing when error
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          app.args,
	}

	cobraComand.SetOut(os.Stdout)
	cobraComand.SetErr(os.Stderr)
	cobraComand.Flags().SortFlags = true

	if len(app.commands) > 0 {
		for _, cmd := range app.commands {
			cobraComand.AddCommand(cmd.cobraCommand())
		}
		cobraComand.SetHelpCommand(helpCommand(app.name))
	}

	if app.runFunc != nil {
		cobraComand.RunE = app.runCommand
	}

	var namedFlagSets flags.NamedFlagSets
	if app.options != nil {
		namedFlagSets = app.options.Flags()
		for _, flagSet := range namedFlagSets.FlagSetMap {
			cobraComand.Flags().AddFlagSet(flagSet)
		}

		usageFmt := "Usage:\n %s\n"
		cols, _, _ := term.TerminalSize(cobraComand.OutOrStdout())
		cobraComand.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
			flags.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
		})
		cobraComand.SetUsageFunc(func(cmd *cobra.Command) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
			flags.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

			return nil
		})
	}

	if !app.noVersion {
		verflag.AddFlags(namedFlagSets.GetFlagSet("global"))
	}

	if !app.noConfig {
		addConfigFlag(app.basename, namedFlagSets.GetFlagSet("global"))
	}

	flags.AddGlobalFlags(namedFlagSets.GetFlagSet("global"), cobraComand.Name())
	app.cmd = &cobraComand
}

func (app *App) Run() {
	if err := app.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func (app *App) Command() *cobra.Command {
	return app.cmd
}

func (app *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	flags.PrintFlags(cmd.Flags())

	if !app.noVersion {
		verflag.PrintAndExitIfRequested()
	}

	if !app.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(app.options); err != nil {
			return err
		}
	}

	if !app.silence {
		log.Infof("%v Starting %s ...", progressMessage, app.name)
		if !app.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !app.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}

	if app.options != nil {
		if err := app.applyOptionRules(); err != nil {
			return err
		}
	}

	if app.runFunc != nil {
		return app.runFunc(app.basename)
	}

	return nil
}

func (app *App) applyOptionRules() error {
	if completeableOptions, ok := app.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := app.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	if printableOptions, ok := app.options.(PrintableOptions); ok && !app.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}
