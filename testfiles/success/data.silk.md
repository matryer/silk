# Echo server

## GET /echo

* `Content-Type`: `"application/json"`
* `X-Another-Header`: `"value"`

```
{"name":"Silk","status":"awesome","a_bool":true,"nothing":null,"release_year":2016}
```

===

### Response

* `Server`: `"EchoDataHandler"`
* `Status`: `200`
* `Data.body.name`: `"Silk"`
* `Data.body.status`: `"awesome"`
* `Data.body.a_bool`: `true`
* `Data.body.nothing`: `null`
* `Data.body.release_year`: `2016`
