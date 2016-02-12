package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/matryer/silk/runner"
	"golang.org/x/build/livelog"
)

type WebRunnerT struct {
	Filepath string
	Fail     bool
	Logger   *log.Logger
	Buffer   *livelog.Buffer
}

func NewWebRunnerT(filepath string) *WebRunnerT {
	buff := &livelog.Buffer{}
	logger := log.New(buff, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &WebRunnerT{
		Filepath: filepath,
		Logger:   logger,
		Buffer:   buff,
	}
}
func (t *WebRunnerT) FailNow() {
	t.Fail = true
}
func (t *WebRunnerT) Log(o ...interface{}) {
	fmt.Println("@@@@@", o)
	t.Logger.Println(o...)
}
func (t *WebRunnerT) LogOutut() string {
	if t == nil {
		return "<empty log>"
	}

	out := string(t.Buffer.Bytes())
	if out == "" {
		return "<empty log>"
	}
	return ""
}

func RunOne(host string, filename string) *WebRunnerT {
	t := NewWebRunnerT(filename)
	runner.New(t, host).RunGlob(filepath.Glob(t.Filepath))
	return t
}

func Run(host string, filepaths ...string) map[string]*WebRunnerT {
	results := map[string]*WebRunnerT{}
	for _, path := range filepaths {
		results[path] = RunOne(host, path)
	}
	return results
}
