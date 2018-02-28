package testutil

import (
	"fmt"
	"strings"
)

// TestT is a test double of type T used to create Runner.
type TestT struct {
	log    []string
	failed bool
}

func (t *TestT) FailNow() {
	t.failed = true
}

func (t *TestT) Failed() bool {
	return t.failed
}

func (t *TestT) LogString() string {
	return strings.Join(t.log, "\n")
}

func (t *TestT) Log(args ...interface{}) {
	t.log = append(t.log, fmt.Sprint(args...))
}
