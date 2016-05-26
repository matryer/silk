package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// EchoHandler gets an http.Handler that echos request data
// back in the response.
func EchoHandler() http.Handler {
	return http.HandlerFunc(handleEcho)
}

// EchoDataHandler gets an http.Handler that echos request data
// back in the response in JSON format.
func EchoDataHandler() http.Handler {
	return http.HandlerFunc(handleEchoData)
}

// EchoRawHandler gets an http.Handler that echos request's body only.
func EchoRawHandler() http.Handler {
	return http.HandlerFunc(handleEchoRaw)
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	// set Server header
	w.Header().Set("Server", "EchoHandler")
	if len(r.Cookies()) > 0 {
		// echo cookies
		for _, cookie := range r.Cookies() {
			http.SetCookie(w, cookie)
		}
	}
	// write summary of request
	fmt.Fprintln(w, strings.ToUpper(r.Method), r.URL.Path)
	// put in the Content-Length
	var bodybuf bytes.Buffer
	if _, err := io.Copy(&bodybuf, r.Body); err != nil {
		log.Println("copying request into buffer failed:", err)
	}
	r.Header.Set("Content-Length", strconv.Itoa(bodybuf.Len()))
	// write parameters
	writeSortedQuery(w, r.URL.Query())
	// write headers
	writeSortedHeaders(w, r.Header)
	// write cookies
	if len(r.Cookies()) > 0 {
		// write sorted cookies out (sorted by name)
		writeSortedCookies(w, r.Cookies(), r)
	}
	// write body
	if _, err := io.Copy(w, &bodybuf); err != nil {
		log.Println("copying request into response failed:", err)
	}
}

func handleEchoData(w http.ResponseWriter, r *http.Request) {
	// set Server header
	w.Header().Set("Server", "EchoDataHandler")

	out := make(map[string]interface{})
	out["method"] = r.Method
	out["path"] = r.URL.Path

	var bodybuf bytes.Buffer
	if _, err := io.Copy(&bodybuf, r.Body); err != nil {
		log.Println("copying request into buffer failed:", err)
	}
	r.Header.Set("Content-Length", strconv.Itoa(bodybuf.Len()))
	for k := range r.Header {
		for _, v := range r.Header[k] {
			out[k] = v
		}
	}
	for k, vs := range r.URL.Query() {
		out[k] = vs
	}
	out["bodystr"] = bodybuf.String()
	var bodyData interface{}
	if err := json.NewDecoder(&bodybuf).Decode(&bodyData); err != nil {
		out["bodyerr"] = err.Error()
	}
	out["body"] = bodyData
	if err := json.NewEncoder(w).Encode(out); err != nil {
		panic(err)
	}
}

func handleEchoRaw(w http.ResponseWriter, r *http.Request) {
	// set Server header
	w.Header().Set("Server", "EchoRawHandler")

	// read body
	var bodybuf bytes.Buffer
	if _, err := io.Copy(&bodybuf, r.Body); err != nil {
		log.Println("copying request into buffer failed:", err)
	}
	// write body
	if _, err := io.Copy(w, &bodybuf); err != nil {
		log.Println("copying request into response failed:", err)
	}
}

func writeSortedHeaders(w io.Writer, headers http.Header) {
	// get header keys
	var keys []string
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range headers[k] {
			vb, err := json.Marshal(v)
			if err != nil {
				log.Println("silk/testutil: cannot marshal header value:", err)
				continue
			}
			fmt.Fprintln(w, "* "+k+":", string(vb))
		}
	}
}

func writeSortedCookies(w io.Writer, cookies []*http.Cookie, r *http.Request) {
	var keys []string
	for _, c := range cookies {
		keys = append(keys, c.Name)
	}
	sort.Strings(keys)
	for _, k := range keys {
		cookie, err := r.Cookie(k)
		if err != nil {
			log.Println("failed to get cookie:", err)
			continue
		}
		fmt.Fprintln(w, "* Cookie:", cookie.String())
	}
}

func writeSortedQuery(w io.Writer, query url.Values) {
	// get header keys
	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		var vals []string
		for _, v := range query[k] {
			vals = append(vals, v)
		}
		sort.Strings(vals)
		for _, v := range vals {
			fmt.Fprintln(w, "* ?"+k+"="+v)
		}
	}
}
