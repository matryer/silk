# Comments

## `POST /comments`

Create a new comment.

* `Content-Type`: `"application/json"`
* `Accept`: `"application/json"`

Include the `name` and `comment` text in the body:

```
{
  "name": "Mat",
  "comment": "Writing tests is easy"
}
```

===

### Example response

* `Status`: `201`
* `Content-Type`: `"application/json"`

```
{
  "id": "123",
  "name": "Mat",
  "comment": "Writing tests is easy"
}
```