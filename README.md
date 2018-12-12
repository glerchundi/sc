# subcommands

Sub-commands package for sub-commands handling. Mainly motivated to combine with [`peterbourgon/ff`](https://github.com/peterbourgon/ff) flag parsing library.

## Usage

```go
func main() {
	cmd := subcommands.NewCommand("myapp", flag.CommandLine, nil)
	
	var addr string
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.StringVar(&addr, "addr", ":8080", "listening addr")
	cmd.AddCommand(subcommands.NewCommand(fs.Name(), fs, func() error {
		return http.ListenAndServe(addr, nil)
	})
	
	err := cmd.Execute(os.Args[1:], func(fs *flag.FlagSet, args []string) error {
		return ff.Parse(fs, args,
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix(strings.ToUpper(cmd.Name())),
		)
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
```
