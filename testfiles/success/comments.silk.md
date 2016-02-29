# Comments and things

* Root: "http://localhost:8080/"

## `POST /comments`

Create a comment.

### Example request

```json
{
  "name":    "Mat",
  "comment": "Good work"
}
```

* `Content-Type`: "application/json" // ensure correct content type is specified

===

### Example response

* `Status`: `201`

```json
{
  "id":      "123",
  "name":    "Mat",
  "comment": "Good work"
}
```

## `GET` `/comments/{id}`

Read a single comment with the specified `{id}`.

* `?pretty=true` // get pretty output

===

* `Status`: `200`
* `Content-Type`: `"application/json"`

```
{
  "id":      "123",
  "name":    "Mat",
  "comment": "Good work"
}
```

# Another group

## DELETE /something/1

===

* `Status`: `200` // OK
