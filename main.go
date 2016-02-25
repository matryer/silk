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

var (
	showVersion    = flag.Bool("version", false, "show version and exit")
	url            = flag.String("silk.url", "", "(required) target url")
	help           = flag.Bool("help", false, "show help")
	root           string
	defaultPattern string
)

func init() {
	root = "."
	defaultPattern = "*.silk.md"
}

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
		root = filepath.Join(root, defaultPattern)
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
	fmt.Println("usage: silk [path/to/files/[pattern]]")
	flag.PrintDefaults()
	fmt.Printf("\nBy default silk will run ./%s\n", defaultPattern)
}

func printversion() {
	fmt.Println("silk", Gitinfo)
}
