# Echo server

## GET /echo

* Content-Type: "text/plain"

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
Hello silk.
```

* Content-Type: "wrong/type"
