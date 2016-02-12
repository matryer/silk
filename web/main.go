package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/unrolled/render"
)

type View struct {
	Files   map[string]File
	Results map[string]*WebRunnerT
}

type File struct {
	Name string
	Path string
	HTML template.HTML
}

func main() {
	silkFolder := "../testfiles/success"
	Server(silkFolder)
}

func Server(folder string) {
	rnd := render.New(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
	})

	files, err := WalkSilkMD(folder)
	results := map[string]*WebRunnerT{}

	if err != nil {
		log.Println("Error parsing md files", err)
	}

	index := func(w http.ResponseWriter, req *http.Request) {
		for k, v := range files {
			t := RunOne("http://localhost:9080", v.Path)
			fmt.Println("@@@@@  Fail - ", t.Fail)
			results[k] = t
		}
		rnd.HTML(w, http.StatusOK, "index", &View{Files: files, Results: results})
	}

	md := func(w http.ResponseWriter, req *http.Request) {
		file := req.URL.Path[len("/files/"):]

		data := template.HTML(files[file].HTML)
		rnd.HTML(w, http.StatusOK, "md", data)
	}

	logs := func(w http.ResponseWriter, req *http.Request) {
		file := req.URL.Path[len("/logs/"):]
		out := results[file].LogOutut()
		rnd.HTML(w, http.StatusOK, "log", out)
	}

	run := func(w http.ResponseWriter, req *http.Request) {
		fail := false
		newRun := map[string]*WebRunnerT{}
		for k, v := range files {
			t := RunOne("http://localhost:9080", v.Path)
			newRun[k] = t
			if t.Fail {
				fail = true
			}
		}
		// FIXME not safe, needs a mutex, or channel
		results = newRun

		rnd.HTML(w, http.StatusOK, "status", fail)
	}

	http.HandleFunc("/files/", md)
	http.HandleFunc("/logs/", logs)

	http.HandleFunc("/run/", run)
	http.HandleFunc("/", index)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	err = http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func WalkSilkMD(folder string) (map[string]File, error) {
	files := map[string]File{}

	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) (e error) {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".silk.md") {
			return nil
		}

		html, err := parseMd(path)
		if err != nil {
			return err
		}

		log.Println("Processed path", path)

		files[f.Name()] = File{
			Name: f.Name(),
			Path: path,
			HTML: template.HTML(html),
		}
		return nil
	})
	return files, err
}

func parseMd(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	all, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	unsafe := blackfriday.MarkdownCommon(all)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return html, nil
}
