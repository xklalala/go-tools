// staticcheck analyses Go code and makes it better.
package main

import (
	"log"
	"os"

	"github.com/xklalala/go-tools/lintcmd"
	"github.com/xklalala/go-tools/simple"
	"github.com/xklalala/go-tools/staticcheck"
	"github.com/xklalala/go-tools/stylecheck"
	"github.com/xklalala/go-tools/unused"
	"golang.org/x/tools/go/analysis"
)

func main() {
	fs := lintcmd.FlagSet("staticcheck")
	debug := fs.String("debug.unused-graph", "", "Write unused's object graph to `file`")
	fs.Parse(os.Args[1:])

	var cs []*analysis.Analyzer
	for _, v := range simple.Analyzers {
		cs = append(cs, v)
	}
	for _, v := range staticcheck.Analyzers {
		cs = append(cs, v)
	}
	for _, v := range stylecheck.Analyzers {
		cs = append(cs, v)
	}

	cs = append(cs, unused.Analyzer)
	if *debug != "" {
		f, err := os.OpenFile(*debug, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		unused.Debug = f
	}

	lintcmd.ProcessFlagSet(cs, fs)
}
