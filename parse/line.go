package parse

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

var (
	commentPrefix = []byte(` //`)
)

// Line represents a single line.
type Line struct {
	Number  int
	Type    LineType
	Bytes   []byte
	Comment []byte
	Regexp  *regexp.Regexp
	detail  *Detail
}

// ParseLine makes a new Line with the given data.
func ParseLine(n int, unsafeText []byte) (*Line, error) {
	linetype := LineTypePlain
	// trim off comments
	var comment []byte
	text := make([]byte, len(unsafeText))
	copy(text, unsafeText)
	if bytes.Contains(text, commentPrefix) {
		segs := bytes.Split(text, commentPrefix)
		text = segs[0]
		comment = segs[1]
	}
	var rx *regexp.Regexp
	for _, item := range matchTypes {
		if regexes[item.R].Match(text) {
			linetype = item.Type
			rx = regexes[item.R]
			break
		}
	}
	// parse the detail now
	var d *Detail
	if linetype == LineTypeDetail || linetype == LineTypeParam {
		var err error
		d, err = parseDetail(text, rx)
		if err != nil {
			return nil, &ErrLine{N: n, Err: err}
		}
	}
	return &Line{
		Number:  n,
		Type:    linetype,
		Bytes:   text,
		Regexp:  rx,
		Comment: comment,
		detail:  d,
	}, nil
}

func (l *Line) String() string {
	return fmt.Sprintf("%d: (%s) %s", l.Number, l.Type, string(l.Bytes))
}

func (l *Line) Detail() *Detail {
	return l.detail
}

var placeholderRegexp = regexp.MustCompile(`{(.*)}`)

// Capture extracts the first placeholder value from the
// comments.
func (l *Line) Capture() string {
	if len(l.Comment) == 0 {
		return ""
	}
	matches := placeholderRegexp.FindSubmatch(l.Comment)
	if len(matches) < 2 {
		return ""
	}
	return string(matches[1])
}

type Lines []*Line

func (l Lines) Bytes() []byte {
	var lines [][]byte
	for _, line := range l {
		lines = append(lines, line.Bytes)
	}
	return bytes.Join(lines, []byte("\n"))
}

func (l Lines) String() string {
	return string(l.Bytes())
}

// Reader makes a new io.Reader that will read the
// bytes from every line.
func (l Lines) Reader() io.Reader {
	return bytes.NewReader(l.Bytes())
}

// Number gets the line number of the first line.
func (l Lines) Number() int {
	if len(l) == 0 {
		return 0
	}
	return l[0].Number
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
	LineTypeParam
)

var lineTypeStrs = map[LineType]string{
	LineTypePlain:        "plain",
	LineTypeGroupHeading: "heading",
	LineTypeRequest:      "request",
	LineTypeCodeBlock:    "codeblock",
	LineTypeDetail:       "detail",
	LineTypeSeparator:    "separator",
	LineTypeParam:        "param",
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
	// * ?param=value
	R:    "^\\s*\\* `?\\?(.*=?.*)`?",
	Type: LineTypeParam,
}, {
	// * Content-Type: application/json
	R:    "^\\s*\\* (.*)",
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
