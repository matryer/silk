package example_test

import (
	"net/http/httptest"
	"testing"

	example "github.com/matryer/silk/examples/failure"
)

func TestHello(t *testing.T) {

	// start test server
	server := httptest.NewServer(example.NewServer())
	defer server.Close()

	// make a new runner using the server URL as the target
	// and run the test file.

	//runner.New(t, server.URL).RunFile("hello.silk.md")

	// NOTE: Uncomment the line above to see a real failure.
	//       Is it commented out so the tests don't fail :)

}
