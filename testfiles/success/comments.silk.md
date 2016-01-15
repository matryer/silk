# Comments and things

* Root: "http://localhost:8080/"

## POST /comments

Create a comment.

### Example request

```
{
  "name":    "Mat",
  "comment": "Good work"
}
```

* Content-Type: "application/json"

===

### Example response

* Status: 201

```
{
  "id":      "123",
  "name":    "Mat",
  "comment": "Good work"
}
```

## GET /comments/{id}

Read a single comment with the specified `{id}`.

===

* Status: 200
* Content-Type: "application/json"

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

* Status: 200
