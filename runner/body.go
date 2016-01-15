package runner

import (
	"encoding/json"
	"io"
)

// ParseJSONBody parses a JSON body.
func ParseJSONBody(r io.Reader) (interface{}, error) {
	var v interface{}
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
