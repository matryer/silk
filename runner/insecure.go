package runner

import (
	"crypto/tls"
	"net/http"
)

// AllowConnectionsToSSLSitesWithoutCerts allows the Runner
// to connect to SSL sites wihtout certification checking.
func (r *Runner) AllowConnectionsToSSLSitesWithoutCerts() {
	r.DoRequest = createInsecureClient().Transport.RoundTrip
}

func createInsecureClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
