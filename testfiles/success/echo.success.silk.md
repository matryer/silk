# Echo server

## `GET /echo`

* `Content-Type`: `"text/plain"`
* `X-Another-Header`: `"value"`

```
Hello silk.
```

===

```
GET /echo
* Accept-Encoding: "gzip"
* Content-Length: "11"
* Content-Type: "text/plain"
* User-Agent: "Go-http-client/1.1"
* X-Another-Header: "value"
Hello silk.
```

* Content-Type: "text/plain; charset=utf-8"
* Server: "EchoHandler"
* Status: 200
* Content-Length: "162"
