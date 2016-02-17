# Hello API

The Hello API just says hello to people, in a very polite way.

## `GET /hello`

Gets a personalised greeting.

  * `?name=Mat` // The name of the person to greet

===

### Example response

* Status: `200`
* Content-Type: `text/plain; charset=utf-8`

Returns the text greeting:

```
Hello Mat.
```
