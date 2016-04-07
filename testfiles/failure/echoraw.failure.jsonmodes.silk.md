# Comments and things

* Root: "http://localhost:8080/"

The server echos the request's body directly.

## `GET /json`

### Example request

```json
{ "name": "Mat", "meta" : { "client" : "tester" } }
```

===

### Example response

* `Status`: `200`

Use the json qualifier in your expected request body to allow differences in order and white space.

```json
{
  "comment": "Good work",
  "name":    "Mat",
  "meta" : {
    "client" : "tester"
  }
}
```
## `GET /json/subset`

Create a comment.

### Example request

```json
{ "name": "Mat",  "comment": "Good work", "meta" : { "api" : 1.0 } }
```

===

### Example response

* `Status`: `200`

Use the `json(mode=subset)` qualifier in your expected request body to only check for
a subset of the response. This allows you to scope your tests or to be more lenient
towards future unrelated changes.

```json(mode=subset)
{
  "name":    "Mat",
  "meta" : { "client" : "tester" }
}
```
