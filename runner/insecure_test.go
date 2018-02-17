// +build integration

package runner

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/silk/testutil"
)

func TestAllowConnectionsToSSLSitesWithoutCerts(t *testing.T) {
	is := is.New(t)
	subT := &testutil.TestT{}
	s := httptest.NewServer(testutil.EchoRawHandler())
	defer s.Close()
	r := New(subT, s.URL)
	r.AllowConnectionsToSSLSitesWithoutCerts()
	req, _ := http.NewRequest("GET", "https://untrusted-root.badssl.com", nil)
	_, err := r.DoRequest(req)
	is.Nil(err)
}
