# Gock

Gock is a simple HTTP mocking server (codes, timeout, random) written in Go.

## Usage

Default response code is `204` (No Content).

Query params can be cumulated (ie. wait + code).

Two modes are available:

* `default`: Gock replies to queries directly
* `proxy`: Gock proxifies queries to another backend

```
go run main.go
```

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

## Configuration

Following environment variables may be set:

| Variable | Description | Default value |
|---|---|---|
| `GOCK_PORT` | HTTP port to listen on | `8000` |
| `GOCK_DEBUG` | Run in debug mode (set to `1`) | `0` |
| `GOCK_MODE` | Gock mode: `default` or `proxy` | `default` |
| `GOCK_PROXY_HOST` | Host to proxify in proxy mode (required) | _none_ |
| `GOCK_PROXY_PORT` | Port to proxify in proxy mode (required) | `80` |
| `GOCK_PROXY_CODE` | Response code in proxy mode | _none_ (backend code) |
| `GOCK_PROXY_WAIT` | Waiting time in proxy mode (in seconds) | `0` |
| `GOCK_PROXY_PERCENT` | Percentage (approximate) on which response code or waiting time apply | `100` |
