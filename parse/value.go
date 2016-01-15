package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type errValue []byte

func (e errValue) Error() string {
	return fmt.Sprintf("invalid value: %s (did you forget quotes?)", string(e))
}

type Value struct {
	Data interface{}
}

func (v *Value) String() string {
	b, err := json.Marshal(v.Data)
	if err != nil {
		panic("silk: cannot marshal value: \"" + fmt.Sprintf("%v", v.Data) + "\": " + err.Error())
	}
	return string(b)
}

func ParseValue(src []byte) *Value {
	var v interface{}
	src = bytes.TrimSpace(src)
	if err := json.Unmarshal(src, &v); err != nil {
		return &Value{Data: string(src)}
	}
	return &Value{Data: v}
}
