package parse_test

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/silk/parse"
)

func TestParser(t *testing.T) {
	is := is.New(t)

	groups, err := parse.ParseFile("../testfiles/success/comments.silk.md", "../testfiles/success/comments2.silk.md")
	is.NoErr(err)

	is.Equal(len(groups), 3)
	is.Equal(groups[0].Filename, "../testfiles/success/comments.silk.md")
	is.Equal(groups[1].Filename, "../testfiles/success/comments.silk.md")
	is.Equal(groups[2].Filename, "../testfiles/success/comments2.silk.md")
	is.Equal(groups[0].Title, "Comments and things")
	is.Equal(groups[1].Title, "Another group")

	group := groups[0]
	is.Equal(len(group.Requests), 2)
	is.Equal(len(group.Details), 1)
	is.Equal(group.Details[0].Detail().Key, "Root")
	is.Equal(group.Details[0].Detail().Value.Data, "http://localhost:8080/")

	req1 := group.Requests[0]
	is.Equal("POST", string(req1.Method))
	is.Equal("/comments", string(req1.Path))
	is.Equal(len(req1.Details), 1)
	is.Equal(req1.Details[0].Detail().Key, "Content-Type")
	is.Equal(req1.Details[0].Detail().Value.Data, "application/json")
	is.Equal(req1.ExpectedDetails[0].Detail().Key, "Status")
	is.Equal(req1.ExpectedDetails[0].Detail().Value.Data, 201)
	is.Equal(req1.Body.String(), `{
  "name":    "Mat",
  "comment": "Good work"
}`)
	is.Equal(req1.ExpectedBody.String(), `{
  "id":      "123",
  "name":    "Mat",
  "comment": "Good work"
}`)

	req2 := group.Requests[1]
	is.Equal("GET", req2.Method)
	is.Equal("/comments/{id}", req2.Path)
	is.Equal(req2.ExpectedDetails[0].Detail().Key, "Status")
	is.Equal(req2.ExpectedDetails[0].Detail().Value.Data, 200)
	is.Equal(req2.ExpectedDetails[1].Detail().Key, "Content-Type")
	is.Equal(req2.ExpectedDetails[1].Detail().Value.Data, "application/json")
	is.Equal(req2.ExpectedBody.String(), `{
  "id":      "123",
  "name":    "Mat",
  "comment": "Good work"
}`)
	is.Equal(req2.ExpectedBody.Number(), 44)

	group = groups[1]
	is.Equal(len(group.Requests), 1)

}
