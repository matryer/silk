package parse

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type errValue []byte

func (e errValue) Error() string {
	return fmt.Sprintf("invalid value: %s (did you forget quotes?)", string(e))
}

func isRegex(v interface{}) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return strings.HasPrefix(s, `/`) && strings.HasSuffix(s, `/`)
}

type Value struct {
	Data interface{}
}

func (v Value) String() string {
	if isRegex(v.Data) {
		return v.Data.(string)
	}
	b, err := json.Marshal(v.Data)
	if err != nil {
		panic("silk: cannot marshal value: \"" + fmt.Sprintf("%v", v.Data) + "\": " + err.Error())
	}
	return string(b)
}

// Equal gets whether the Data and specified value are equal.
// Supports regexp values.
func (v Value) Equal(val interface{}) bool {
	var str string
	var ok bool
	if str, ok = v.Data.(string); !ok {
		return v.Data == val
	}
	if isRegex(str) {
		// looks like regexp to me
		regex := regexp.MustCompile(str[1 : len(str)-1])
		// turn the value into a string
		valStr := fmt.Sprintf("%v", val)
		if regex.Match([]byte(valStr)) {
			return true
		}
	}
	return fmt.Sprintf("%v", v.Data) == fmt.Sprintf("%v", val)
}

func (v Value) Type() string {
	var str string
	var ok bool
	if str, ok = v.Data.(string); !ok {
		return fmt.Sprintf("%T", v.Data)
	}
	if isRegex(str) {
		return "regex"
	}
	return "string"
}

func ParseValue(src []byte) *Value {
	var v interface{}
	src = clean(src)
	if err := json.Unmarshal(src, &v); err != nil {
		return &Value{Data: string(src)}
	}
	return &Value{Data: v}
}
