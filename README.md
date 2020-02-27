# Gock

Gock is a simple HTTP mocking server (codes, timeout, random) written in Go.

## Usage

Default response code is `204` (No Content).

Query params can be cumulated (ie. wait + code).

### Request with wait (for timeout tests)

Wait given amount of seconds before response:

```
http://gock:8000/?wait=10
```

### Request with given response code

Respond with `404` HTTP code:

```
http://gock:8000/?code=404
```
