# Echo server

## `GET /echo`

* `Content-Type`: `"text/plain"`
* `X-Another-Header`: `"value"`
* `?param1=value1`
* `?param2=value2`
* `?param3=value3`
* Cookie: name=silk; another=true

```
Hello silk.
```

===

```
GET /echo
* ?param1=value1
* ?param2=value2
* ?param3=value3
* Accept-Encoding: "gzip"
* Content-Length: "11"
* Content-Type: "text/plain"
* Cookie: "name=silk; another=true"
* User-Agent: "Go-http-client/1.1"
* X-Another-Header: "value"
* Cookie: another=true
* Cookie: name=silk
Hello silk.
```

* Content-Type: "text/plain; charset=utf-8"
* Server: "EchoHandler"
* Status: 200
* Content-Length: "292"
* Set-Cookie: /another=true/
* Set-Cookie: /name=silk/
