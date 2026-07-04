// Command yup-cat is the CLI wrapper around github.com/gloo-foo/cmd-cat.
package main

import (
	clix "github.com/gloo-foo/cli"
	command "github.com/gloo-foo/cmd-cat"
	urf "github.com/urfave/cli/v3"
)

// version is the build version. It defaults to "dev" for local builds and is
// overridden at release time via the linker: -ldflags "-X main.version=<v>".
var version = "dev"

const (
	name               = "cat"
	flagNumber         = "number"
	flagNumberNonBlank = "number-nonblank"
)

// synopsis is the multi-line --help usage block. urfave/cli indents the whole
// block three spaces, so the lines stay flush-left.
const synopsis = `cat [OPTIONS] [FILE...]

Concatenate FILE(s) to standard output.
With no FILE, or when FILE is -, read standard input.`

// spec declares the cat wrapper: a file-or-stdin filter with -n/-b numbering.
var spec = clix.Spec{
	Name:     name,
	Summary:  "concatenate files and print on the standard output",
	Synopsis: synopsis,
	Build:    build,
	Flags: []urf.Flag{
		&urf.BoolFlag{Name: flagNumber, Aliases: []string{"n"}, Usage: "number all output lines"},
		&urf.BoolFlag{
			Name:    flagNumberNonBlank,
			Aliases: []string{"b"},
			Usage:   "number nonempty output lines (overrides -n)",
		},
	},
}

// build maps the invocation to cat's pipeline: a file-or-stdin source into the
// cat command configured by the numbering flags.
func build(inv clix.Invocation) (clix.Source, clix.Command, error) {
	return clix.OperandsOrStdin(inv), command.Cat(options(inv.Args)...), nil
}

// options folds the parsed flags into cat's option values.
func options(c *urf.Command) []any {
	var opts []any
	if c.Bool(flagNumber) {
		opts = append(opts, command.CatNumberLines)
	}
	if c.Bool(flagNumberNonBlank) {
		opts = append(opts, command.CatNumberNonBlank)
	}
	return opts
}

// runMain is an indirection seam so main's wiring is testable without spawning
// the process; a test swaps it and restores it.
var runMain = clix.Main

func main() { runMain(spec, version) }
