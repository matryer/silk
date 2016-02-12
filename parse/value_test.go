package parse

import (
	"encoding/json"
	"testing"

	"github.com/cheekybits/is"
)

func TestValue(t *testing.T) {
	is := is.New(t)

	var tests = []interface{}{
		"String",
		123,
		1.23,
		true,
		nil,
	}
	for _, test := range tests {
		b, err := json.Marshal(test)
		is.NoErr(err)
		actual := ParseValue(b)
		is.Equal(actual.Data, test)
	}

}

func TestValueEqual(t *testing.T) {
	is := is.New(t)

	v := ParseValue([]byte("something"))
	is.True(v.Equal("something"))
	is.False(v.Equal("else"))
	is.Equal("string", v.Type())

	v = ParseValue([]byte("123"))
	is.Equal("float64", v.Type())

	v = ParseValue([]byte("/^2.{2}$/"))
	is.True(v.Equal(200))
	is.True(v.Equal(201))
	is.False(v.Equal(404))
	is.Equal("regex", v.Type())

	v = ParseValue([]byte("/application/json/"))
	is.True(v.Equal("application/json"))
	is.True(v.Equal("application/json; charset=utf-8"))
	is.True(v.Equal("text/xml; application/json; charset=utf-8"))
	is.False(v.Equal("text/xml; charset=utf-8"))
	is.Equal("regex", v.Type())
	is.Equal(`/application/json/`, v.String())

	v = ParseValue([]byte("/Silk/"))
	is.True(v.Equal("My name is Silk."))
	is.True(v.Equal("Silk is my name."))
	is.False(v.Equal("I don't contain that word!"))
}
