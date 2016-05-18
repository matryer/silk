# Comments and things

* Root: "http://localhost:8080/"

The server echos the request's body directly.

## `GET /json`

### Example request

```json
{ "name": "Mat",  "comment": "Good work", "meta" : { "api" : 1.0 } }
```

===

### Example response

* `Status`: `200`

By defaul using the `json` qualifier in your expected request body only checks for
a subset of the response. This allows you to scope your tests or to be more lenient
towards future unrelated changes. Additional fields in the response do not invalidate
the test.

```json
{
  "name":    "Mat",
  "meta" : { "client" : "tester" }
}
```

## `GET /json/same`

### Example request

```json
{ "name": "Mat", "meta" : { "client" : "tester" } }
```

===

### Example response

* `Status`: `200`

Use the `json(strict)` qualifier in your expected request body to ensure that the json object are the same.
This allows differences in order and white space.

```json(strict)
{
  "comment": "Good work",
  "name":    "Mat",
  "meta" : {
    "client" : "tester"
  }
}
```
