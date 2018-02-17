package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/juacompe/silk/runner"
)

var (
	showVersion = flag.Bool("version", false, "show version and exit")
	insecure    = flag.Bool("insecure", false, "allow connections to SSL sites without certs")
	url         = flag.String("silk.url", "", "(required) target url")
	help        = flag.Bool("help", false, "show help")
	paths       []string
)

func main() {
	flag.Parse()
	if *showVersion {
		printversion()
		return
	}
	if *help {
		printhelp()
		return
	}
	if *url == "" {
		fmt.Println("silk.url argument is required")
		return
	}
	paths = flag.Args()
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{{Name: "silk", F: testFunc}},
		nil,
		nil)
}

func testFunc(t *testing.T) {
	r := runner.New(t, *url)
	if *insecure {
		r.AllowConnectionsToSSLSitesWithoutCerts()
	}
	fmt.Println("silk: running", len(paths), "file(s)...")
	r.RunGlob(paths, nil)
}

func printhelp() {
	printversion()
	fmt.Println("usage: silk [file] [file2 [file3 [...]]")
	fmt.Println("  e.g: silk ./test/*.silk.md")
	flag.PrintDefaults()
}

func printversion() {
	fmt.Println("silk", version)
}
