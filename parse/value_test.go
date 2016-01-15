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
