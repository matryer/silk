# Echo server

## GET /echo

* `Content-Type`: `"application/json"`
* `X-Another-Header`: `"value"`

```
{"name":"{$AppNameFromEnv}","status":"awesome","a_bool":true,"nothing":null,"release_year":2016,"exp":200}
```

===

### Response

* `Server`: `"EchoDataHandler"`
* `Status`: `200` // Expected {value}
* `Data.body.name`: `"Silk"`
* `Data.body.status`: /awesome/ // The {status} of Silk.
* `Data.body.a_bool`: `true`
* `Data.body.nothing`: `null`
* `Data.body.release_year`: `2016`
* `Data.body.exp`: `200`

## POST /echo/{status}

* ?status={status}
* X-Status: {status}

```
{"st":"{status}"}
```

===

### Response

* `Server`: `"EchoDataHandler"`
* `Status`: `{value}`
* `Data.body.st`: `"awesome"`
* `Data.body.st`: {status}

```
{"Accept-Encoding":"gzip","Content-Length":"16","User-Agent":"Go-http-client/1.1","X-Status":"awesome","body":{"st":"awesome"},"bodystr":"{\"st\":\"awesome\"}","method":"POST","path":"/echo/awesome","status":["awesome"]}

```

## POST /echo/{$EnvStatus}

* ?status={$EnvStatus}
* X-Status: {$EnvStatus}

```
{"st":"{$EnvStatus}"}
```

===

### Response

* `Server`: `"EchoDataHandler"`
* `Status`: `{value}`
* `Data.body.st`: `"awesome"`
* `Data.body.st`: {$EnvStatus}

```
{"Accept-Encoding":"gzip","Content-Length":"16","User-Agent":"Go-http-client/1.1","X-Status":"awesome","body":{"st":"awesome"},"bodystr":"{\"st\":\"awesome\"}","method":"POST","path":"/echo/awesome","status":["awesome"]}

```