![Silk logo](https://github.com/matryer/silk/blob/master/other/SilkLogo-256.png)

# silk

Markdown based document-driven web API testing.

  * Write nice looking Markdown documentation ([like this](https://github.com/matryer/silk/blob/master/testfiles/success/example.silk.md)), and then run it using the [silk command](#command-line)
  * Simple and robust [Markdown API](#markdown-api)
  * Comes with [real examples](https://github.com/matryer/silk/tree/master/testfiles/success) that you can copy (that are also part of the test suite for the project)

PROJECT STATUS: Looking for BETA testers.

![Example Silk test file](https://github.com/matryer/silk/blob/master/other/example.png)

## Markdown API

Tests are made up of documents written in Markdown.

  * `# Group` - Top level headings represent groups of requests
  * `## GET /path` -  Second level headings represent a request
  * Code blocks with three back tics represent bodies
  * `* Field: value` - Lists describe headers and assertions
  * `* ?param=value` - Request parameters
  * `===` seperators break requests from responses
  * Comments (starting with `//`) are ignored
  * Plain text is ignored to allow you to add documentation
  * Inline back tics are ignored and are available for formatting

### Document structure

A document is made up of:

  * A request
  * `===` seperator
  * Assertions

### Requests

A request starts with `##` and must have an HTTP method, and a path:

```
## METHOD /path
```

Examples include:

```
## GET /people

## POST /people/1/comments

## DELETE /people/1/comments/2

```

#### Request body (optional)

To specify a request body (for example for `POST` requests) use a codeblock using backtics (` ``` `):

    ```
    {"name": "Silk", "release_year": 2016}
    ```

#### Request headers (optional)

You may specify request headers using lists (prefixed with `*`):

```
* Content-Type: "application/json"
* X-Custom-Header: "123"
```

#### Request parameters (optional)

Adding parameters to the path (like `GET /path?q=something`) can be tricky, especially when you consider escaping etc. To address this, Silk supports parameters like lists:

```
* ?param=value
```

The parameters will be correctly added to the URL path before the request is made.

### Assertions

Following the `===` separator, you can specify assertions about the response. At a minimum, it is recommended that you assert the status code to ensure the request succeeded:

```
  * Status: 200
```

You may also specify response headers in the same format as request headers:

```
  * Content-Type: "application/json"
  * X-MyServer-Version: "v1.0"
```

If any of the headers do not match, the test will fail.

#### Validating data

You can optionally include a verbatim body using ` ``` ` code blocks. If the response body does not exactly match, the test will fail:

    ```
    {"id": 1, "name": "Silk", "release_year": 2016}
    ```

Alternatively, you can specify a list (using `*`) of data fields to assert accessible via the `Data` object:

```
  * Status: 201
  * Content-Type: "application/json"
  * Data.name: "Silk"
  * Data.release_year: 2016
  * Data.tags[0]: "testing"
  * Data.tags[1]: "markdown"
```

  * NOTE: Currenly this feature is only supported for JSON APIs.

#### Regex

Values may be regex, if they begin and end with a forward slash: `/`. The assertion will pass if the value (after being turned into a string) matches the regex.

```
  * Status: /^2.{2}$/
  * Content-Type: /application/json/
```

The above will assert that:

  * The status beings with `2`, and
  * The `Content-Type` contains `application/json`

## Command line

The `silk` command runs tests against an HTTP endpoint.

Usage:

```
silk -silk.url="{endpoint}" {testfiles}
```

  * `{url}` the endpoint URL (e.g. `http://localhost:8080`)
  * `{testfiles}` path to test files (e.g. `./testfiles`)

Notes:

  * Omit trailing slash from `url`
  * `{testfiles}` can include a pattern (e.g. `/path/*.silk.md`)

## Golang

Silk is written in Go and integrates seamlessly into existing testing tools and frameworks. Import the `runner` package and use `RunGlob` to match many test files:

```
package project_test

import (
  "testing"
  "github.com/matryer/silk/runner"
)

func TestAPIEndpoint(t *testing.T) {
  // start a server
  s := httptest.NewServer(yourHandler)
  defer s.Close()
  
  // run all test files
  runner.New(t, s.URL).RunGlob(filepath.Glob("../testfiles/failure/*.silk.md"))
}
```

  * See the [documentation for the silk/runner package](https://godoc.org/github.com/matryer/silk/runner)
