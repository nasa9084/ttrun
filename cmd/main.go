package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	flags "github.com/jessevdk/go-flags"
	"github.com/nasa9084/errors"
)

type options struct {
	Args    []string `short:"a" long:"arg" description:"argument formed foo=bar"`
	Verbose bool     `short:"v" long:"verbose"`
	PosArgs struct {
		Template string `positional-arg-name:"TEMPLATE" description:"template filepath"`
	} `positional-args:"yes" required:"yes"`
}

func main() {
	if err := _main(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func _main() error {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		if fe, ok := err.(*flags.Error); ok && fe.Type == flags.ErrHelp {
			return nil
		}
		return errors.Wrap(err, "parsing options")
	}
	args, err := parseArgs(opts.Args)
	if err != nil {
		return err
	}
	if opts.Verbose {
		printArgs(args)
	}
	t, err := template.ParseFiles(opts.PosArgs.Template)
	if err != nil {
		return errors.Wrap(err, "parsing template file")
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, args); err != nil {
		return errors.Wrap(err, "executing template")
	}
	if _, err := buf.WriteTo(os.Stdout); err != nil {
		return errors.Wrap(err, "writing to stdout")
	}
	return nil
}

func parseArgs(args []string) (map[string]string, error) {
	ret := map[string]string{}
	for _, arg := range args {
		if !strings.Contains(arg, "=") {
			key := strings.Trim(arg, " ")
			if _, ok := ret[key]; ok {
				return nil, fmt.Errorf("argument %s is already defined", key)
			}
			ret[key] = "true"
			continue
		}
		fs := strings.FieldsFunc(arg, func(r rune) bool { return r == '=' })
		if len(fs) != 2 {
			return nil, fmt.Errorf("argument should be formed foo=bar")
		}
		key := strings.Trim(fs[0], " ")
		val := strings.Trim(fs[1], " ")
		if _, ok := ret[key]; ok {
			return nil, fmt.Errorf("argument %s is already defined", key)
		}
		ret[key] = val
	}
	return ret, nil
}

func printArgs(args map[string]string) {
	for k, v := range args {
		fmt.Printf("%s: %v\n", k, v)
	}
	fmt.Println("==========")
}
