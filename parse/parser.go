package parse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	errMissingGroupHeader  = errors.New("missing group header")
	errUnexpectedCodeblock = errors.New("unexpected codeblock")
	errMissingEndCodeblock = errors.New("missing end codeblock")
	errUnexpectedDetails   = errors.New("unexpected details")
	errUnexpectedParams    = errors.New("unexpected params")
	errMalformedDetail     = errors.New("malformed detail")
)

// Group represents a group of Requests.
type Group struct {
	Filename string
	Title    []byte
	Requests []*Request
	Details  Lines
}

// Request describes an HTTP request and a set of
// associated assertions.
type Request struct {
	Path     []byte
	Method   []byte
	Details  Lines
	Params   Lines
	Body     Lines
	BodyType string
	//===
	ExpectedBody     Lines
	ExpectedBodyType string
	ExpectedDetails  Lines
}

// ErrLine describes an error at a specific line.
type ErrLine struct {
	N   int
	Err error
}

func (e ErrLine) Error() string {
	return fmt.Sprintf("%d: %v", e.N, e.Err)
}

// ParseFile parses the specified files.
func ParseFile(files ...string) ([]*Group, error) {
	var groups []*Group
	for _, file := range files {
		if err := func(file string) error {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			gs, err := Parse(file, f)
			if err != nil {
				return err
			}
			groups = append(groups, gs...)
			return nil
		}(file); err != nil {
			return nil, err
		}
	}
	return groups, nil
}

// Parse parses a file.
func Parse(filename string, r io.Reader) ([]*Group, error) {

	n := 0
	groups := make([]*Group, 0)
	scanner := bufio.NewScanner(r)

	// whether we're at the point of expectations or
	// not.
	settingExpectations := false

	var currentGroup *Group
	var currentRequest *Request

	for scanner.Scan() {
		n++
		line, err := ParseLine(n, scanner.Bytes())
		if err != nil {
			return nil, err
		}
		switch line.Type {
		case LineTypeGroupHeading:
			// new group
			if currentGroup != nil {
				if currentRequest != nil {
					currentGroup.Requests = append(currentGroup.Requests, currentRequest)
					currentRequest = nil
				}
				groups = append(groups, currentGroup)
			}
			title, err := getok(line.Regexp.FindSubmatch(line.Bytes), 1)
			if err != nil {
				return nil, &ErrLine{N: n, Err: err}
			}
			currentGroup = &Group{
				Filename: filename,
				Title:    title,
			}
		case LineTypeRequest:
			// new request
			if currentGroup == nil {
				return nil, &ErrLine{N: n, Err: errMissingGroupHeader}
			}
			if currentRequest != nil {
				currentGroup.Requests = append(currentGroup.Requests, currentRequest)
			}
			settingExpectations = false
			var err error
			currentRequest = &Request{}
			matches := line.Regexp.FindSubmatch(line.Bytes)
			if currentRequest.Method, err = getok(matches, 1); err != nil {
				return nil, &ErrLine{N: n, Err: err}
			}
			if currentRequest.Path, err = getok(matches, 2); err != nil {
				return nil, &ErrLine{N: n, Err: err}
			}
		case LineTypeCodeBlock:

			if currentRequest == nil {
				return nil, &ErrLine{N: n, Err: errUnexpectedCodeblock}
			}

			var bodyType string
			if len(line.Bytes) > 3 {
				bodyType = string(line.Bytes[3:])
			}

			var lines Lines
			var err error
			n, lines, err = scancodeblock(n, scanner)
			if err != nil {
				return nil, &ErrLine{N: n, Err: err}
			}
			if settingExpectations {
				currentRequest.ExpectedBody = lines
				currentRequest.ExpectedBodyType = bodyType
			} else {
				currentRequest.Body = lines
				currentRequest.BodyType = bodyType
			}

		case LineTypeDetail:
			if currentRequest == nil && currentGroup == nil {
				return nil, &ErrLine{N: n, Err: errUnexpectedDetails}
			}
			if currentRequest == nil {
				currentGroup.Details = append(currentGroup.Details, line)
				continue
			}
			if settingExpectations {
				currentRequest.ExpectedDetails = append(currentRequest.ExpectedDetails, line)
			} else {
				currentRequest.Details = append(currentRequest.Details, line)
			}
		case LineTypeParam:
			if currentRequest == nil && currentGroup == nil {
				return nil, &ErrLine{N: n, Err: errUnexpectedParams}
			}
			if settingExpectations {
				return nil, &ErrLine{N: n, Err: errUnexpectedParams}
			}
			currentRequest.Params = append(currentRequest.Params, line)
		case LineTypeSeparator:
			settingExpectations = true
		}

	}

	if currentGroup == nil {
		return nil, &ErrLine{N: n, Err: errMissingGroupHeader}
	}
	if currentRequest != nil {
		currentGroup.Requests = append(currentGroup.Requests, currentRequest)
	}
	groups = append(groups, currentGroup)

	return groups, nil
}

func scancodeblock(n int, scanner *bufio.Scanner) (int, Lines, error) {
	var lines Lines
	for scanner.Scan() {
		n++
		line, err := ParseLine(n, scanner.Bytes())
		if err != nil {
			return n, nil, err
		}
		if line.Type == LineTypeCodeBlock {
			// we're done
			return n, lines, nil
		}
		lines = append(lines, line)
	}
	// shouldn't reach the end
	return n, lines, errMissingEndCodeblock
}

func getok(src [][]byte, i int) ([]byte, error) {
	if i+1 > len(src) {
		return nil, fmt.Errorf("bad format: expected at least %d regex matches, but was %d: %s", i+1, len(src), string(bytes.Join(src, []byte("\n"))))
	}
	return clean(src[i]), nil
}
