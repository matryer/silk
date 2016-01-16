package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/silk/runner"
)

/*
	Usage:

		silk [path] -p=../*.silk.md
*/

var (
	showVersion = flag.Bool("version", false, "show version and exit")
	url         = flag.String("silk.url", "", "(required) target url")
	help        = flag.Bool("help", false, "show help")
	root        string
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
	if len(*url) == 0 {
		fmt.Println("must provide -silk.url")
		printhelp()
		return
	}
	root = "."
	args := flag.Args()
	if len(args) > 0 {
		root = args[0]
	}
	info, err := os.Stat(root)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if info.IsDir() {
		// add default pattern
		root = filepath.Join(root, "*.silk.md")
	}
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{{Name: "silk", F: testFunc}},
		nil,
		nil)
}

func testFunc(t *testing.T) {
	r := runner.New(t, *url)
	files, err := filepath.Glob(root)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("running", len(files), "file(s)")
	r.RunGlob(files, nil)
}

func printhelp() {
	printversion()
	fmt.Println(`usage:
  silk [path/to/files/[pattern]]`)
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println(`By default silk will run ./*.silk.md`)
}

func printversion() {
	fmt.Println("silk", version)
}
