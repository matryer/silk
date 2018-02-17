// +build integration

package runner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/silk/testutil"
)

func TestAllowConnectionsToSSLSitesWithoutCerts(t *testing.T) {
	is := is.New(t)
	subT := &testT{}
	s := httptest.NewServer(testutil.EchoRawHandler())
	defer s.Close()
	r := New(subT, s.URL)
	r.AllowConnectionsToSSLSitesWithoutCerts()
	req, _ := http.NewRequest("GET", "https://untrusted-root.badssl.com", nil)
	_, err := r.DoRequest(req)
	is.Nil(err)
}

type testT struct {
	log    []string
	failed bool
}

func (t *testT) FailNow() {
	t.failed = true
}

func (t *testT) Failed() bool {
	return t.failed
}

func (t *testT) LogString() string {
	return strings.Join(t.log, "\n")
}

func (t *testT) Log(args ...interface{}) {
	t.log = append(t.log, fmt.Sprint(args...))
}
