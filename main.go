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
	url  = flag.String("url", "", "target url")
	help = flag.Bool("help", false, "show help")
	root string
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println(`usage:
  silk [path/to/files/[pattern]]`)
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println(`By default silk will match *.silk.md files in the current directory.`)
		return
	}
	if len(*url) == 0 {
		fmt.Println("must provide -url")
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
		[]testing.InternalTest{{Name: "silk", F: all}},
		nil,
		nil)
}

func all(t *testing.T) {
	r := runner.New(t, *url)
	files, err := filepath.Glob(root)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("running", len(files), "file(s)")
	r.RunGlob(files, nil)
}
