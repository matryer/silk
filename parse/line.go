package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var (
	commentPrefix = []byte(` //`)
)

// Line represents a single line.
type Line struct {
	Number int
	Type   LineType
	Bytes  []byte
	Regexp *regexp.Regexp
	detail *Detail
}

// NewLine makes a new Line with the given data.
func NewLine(n int, text []byte) (*Line, error) {
	linetype := LineTypePlain
	var rx *regexp.Regexp
	for _, item := range matchTypes {
		if regexes[item.R].Match(text) {
			linetype = item.Type
			rx = regexes[item.R]
			break
		}
	}
	// trim off comments
	if bytes.Contains(text, commentPrefix) {
		text = bytes.Split(text, commentPrefix)[0]
	}
	// parse the detail now
	var d *Detail
	if linetype == LineTypeDetail {
		var err error
		d, err = parseDetail(text, rx)
		if err != nil {
			return nil, &ErrLine{N: n, Err: err}
		}
	}
	return &Line{
		Number: n,
		Type:   linetype,
		Bytes:  text,
		Regexp: rx,
		detail: d,
	}, nil
}

func (l *Line) String() string {
	return fmt.Sprintf("%d: (%s) %s", l.Number, l.Type, string(l.Bytes))
}

func (l *Line) Detail() *Detail {
	return l.detail
}

type Lines []*Line

func (l Lines) Join() []byte {
	var lines [][]byte
	for _, line := range l {
		lines = append(lines, line.Bytes)
	}
	return bytes.Join(lines, []byte("\n"))
}

func (l Lines) String() string {
	return string(l.Join())
}

// Reader gets an io.Reader that will read every line.
func (l Lines) Reader() io.Reader {
	var readers []io.Reader
	for _, line := range l {
		readers = append(readers, bytes.NewReader(line.Bytes))
	}
	return io.MultiReader(readers...)
}

// Number gets the line number of the first line.
func (l Lines) Number() int {
	if len(l) == 0 {
		return 0
	}
	return l[0].Number
}

type Detail struct {
	Key   string
	Value *Value
}

func parseDetail(b []byte, detailregex *regexp.Regexp) (*Detail, error) {
	detail, err := getok(detailregex.FindSubmatch(b), 1)
	if err != nil {
		panic("silk: failed to parse detail: " + err.Error())
	}
	colon := bytes.Index(detail, []byte(":"))
	if colon == -1 || colon > len(detail)-1 {
		return nil, errors.New("malformed detail")
	}
	key := clean(detail[0:colon])
	return &Detail{
		Key:   string(bytes.TrimSpace(key)),
		Value: ParseValue(detail[colon+1:]),
	}, nil
}

func (d *Detail) String() string {
	valbytes, err := json.Marshal(d.Value)
	if err != nil {
		return d.Key + ": " + fmt.Sprint(d.Value)
	}
	return d.Key + ": " + string(valbytes)
}

func clean(b []byte) []byte {
	return bytes.Trim(bytes.TrimSpace(b), "`")
}

// LineType represents the type of a line.
type LineType int8

// LineTypes
const (
	LineTypePlain LineType = iota
	LineTypeGroupHeading
	LineTypeRequest
	LineTypeCodeBlock
	LineTypeDetail
	LineTypeSeparator
)

var lineTypeStrs = map[LineType]string{
	LineTypePlain:        "plain",
	LineTypeGroupHeading: "heading",
	LineTypeRequest:      "request",
	LineTypeCodeBlock:    "codeblock",
	LineTypeDetail:       "detail",
	LineTypeSeparator:    "separator",
}

func (l LineType) String() string {
	return lineTypeStrs[l]
}

// matchTypes map patterns to types.
// Prescedence is important.
var matchTypes = []struct {
	R    string
	Type LineType
}{{
	// ## GET /comments
	R:    "^## (.*) (.*)",
	Type: LineTypeRequest,
}, {
	// # Heading
	R:    "^# (.*)",
	Type: LineTypeGroupHeading,
}, {
	// ```
	R:    "^```",
	Type: LineTypeCodeBlock,
}, {
	// ===
	R:    "^(===+)",
	Type: LineTypeSeparator,
}, {
	// ---
	R:    "^(---+)",
	Type: LineTypeSeparator,
}, {
	// - Content-Type: application/json
	R:    "^- (.*)",
	Type: LineTypeDetail,
}, {
	// * Content-Type: application/json
	R:    "^\\* (.*)",
	Type: LineTypeDetail,
}}

var regexes map[string]*regexp.Regexp

func init() {
	// compile regexes
	regexes = make(map[string]*regexp.Regexp, len(matchTypes))
	for _, item := range matchTypes {
		regexes[item.R] = regexp.MustCompile(item.R)
	}
}
