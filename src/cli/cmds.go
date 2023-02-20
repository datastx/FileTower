package cli

// CLI options
type CLI struct {
	Directory string `help:"Directory to monitor"`
	Config    string `help:"Location of configuration file"`
}
