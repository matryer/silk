package main

import (
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

var (
	PageView *View
	rnd      *render.Render
)

type View struct {
	Host       string
	SilkFolder string

	Files   map[string]File
	Results map[string]*WebRunnerT

	Fail bool
}

type File struct {
	Name string
	Path string
	HTML template.HTML
}

// Run the silk test against host for example "http://localhost:9080"
func (v *View) Run(host string) {
	fail := false
	results := map[string]*WebRunnerT{}
	for k, v := range v.Files {
		t := RunOne(host, v.Path)
		results[k] = t
		if t.Fail {
			fail = true
		}
	}
	// FIXME not safe, needs a mutex, or channel
	v.Fail = fail
	v.Results = results
}

func (v *View) LoadFiles(folder string) {
	v.SilkFolder = folder
	files, err := WalkSilkMD(folder)
	if err != nil {
		log.Println("Error parsing md files", err)
	}

	v.Files = files
}

func main() {
	silkFolder := "../testfiles/success"
	Server(silkFolder)
}

func Server(folder string) {
	rnd = render.New(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
	})

	// global state
	PageView = &View{}
	PageView.Host = "http://localhost:9080"
	PageView.LoadFiles(folder)

	http.HandleFunc("/files/", MarkdownHandler)
	http.HandleFunc("/logs/", LogsHandler)

	http.HandleFunc("/run/", RunHandler)
	http.HandleFunc("/", IndexHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	PageView.Run(PageView.Host)
	rnd.HTML(w, http.StatusOK, "index", PageView)
}

func MarkdownHandler(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Path[len("/files/"):]

	data := template.HTML(PageView.Files[file].HTML)
	rnd.HTML(w, http.StatusOK, "md", data)
}

func LogsHandler(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Path[len("/logs/"):]
	out := PageView.Results[file].LogOutut()
	rnd.HTML(w, http.StatusOK, "log", out)
}

func RunHandler(w http.ResponseWriter, req *http.Request) {
	host := req.FormValue("host")
	PageView.Run(host)
	rnd.HTML(w, http.StatusOK, "navstatus", PageView)
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
