package subcommands

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Command represents a single command.
type Command struct {
	name        string
	flagSet     *flag.FlagSet
	run         func() error
	subcommands map[string]*Command
}

// NewCommand returns a new command with the specified name, flags and runner.
// The runner could be nil in order to use the command as the root of other
// subcommands.
func NewCommand(name string, flagSet *flag.FlagSet, run func() error) *Command {
	c := &Command{
		name:    name,
		flagSet: flagSet,
		run:     run,
	}
	c.flagSet.Usage = func() { c.usage() }
	return c
}

// Name returns the name of the command.
func (c *Command) Name() string {
	return c.name
}

// Execute executes the command and returns an error in case the runner fails.
func (c *Command) Execute(args []string, p FlagParser) error {
	c, err := c.traverse(args, p)
	if err != nil {
		return err
	}

	if c.run == nil {
		c.usage()
		os.Exit(2)
	}

	return c.run()
}

// AddCommand adds a subcommand to the supported subcommands.
func (c *Command) AddCommand(subc *Command) {
	if c.subcommands == nil {
		c.subcommands = make(map[string]*Command)
	}

	c.subcommands[subc.name] = subc
}

func (c *Command) traverse(args []string, p FlagParser) (*Command, error) {
	if err := p(c.flagSet, args); err != nil {
		return nil, err
	}

	name := c.flagSet.Arg(0)
	subc, ok := c.subcommands[name]
	if !ok {
		return c, nil
	}

	return subc.traverse(c.flagSet.Args()[1:], p)
}

func (c *Command) usage() {
	fs := c.flagSet

	fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
	fs.PrintDefaults()

	if len(c.subcommands) > 0 {
		names := make([]string, 0, len(c.subcommands))
		for name := range c.subcommands {
			names = append(names, name)
		}
		fmt.Fprintf(fs.Output(), "\nAvailable commands: %s.\n", strings.Join(names, ", "))
	}
}

// FlagParser parses a flag set with the specified arguments.
type FlagParser func(fs *flag.FlagSet, args []string) error
