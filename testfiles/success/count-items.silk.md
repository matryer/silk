# Echo server


## GET /echo/array_fields

Count items in array.

```
{"numbers": [1,2,3,4,5]}
```
===

### Response

* `Server`: `"EchoDataHandler"`
* `Status`: `200`
* `Data.body.numbers.length`: 5


## GET /echo/strings

Count characters in strings.

* `Content-Type`: `"application/json"`
* `X-Another-Header`: `"value"`

```
{"ascii": "Hello World.", "utf-8": "こんにちは、世界。"}

```

===

### Response

* `Server`: `"EchoDataHandler"`
* `Status`: `200`
* `Data.body.ascii.length`: 12
* `Data.body.utf-8.length`: 9

