# coraza-httpbin

`httpbin` server with coraza as reverse proxy to play with attacks.

## Getting started

### Running in Go

First you need to run the server using a directives file:

```shell
go run github.com/jcchavezs/coraza-httpbin/cmd/coraza-httpbin@latest -directives ./directives.conf.example
```

and do a HTTP call in another terminal:

```shell
# attempt XSS injection
curl -i "localhost:8080?arg=<script>a</script>"
```

**Important:** Notice `@owasp_crs` folder is already included and can be used as described [here](https://github.com/corazawaf/coraza-coreruleset).

You can also point to rules in your filesystem by using the absolute notation in the directives file:

```seclang
SecRuleEngine On
SecDebugLog /dev/stdout
SecDebugLogLevel 1
Include /path-to-my-rules/coraza.conf-recommended
Include /path-to-my-rules/crs-setup.conf.example
Include /path-to-my-rules/*.conf
```

### Running in Docker

```shell
docker run -v /path-to-rules:/etc/my-rules ghcr.io/jcchavezs/coraza-httpbin:main
```

and change your directives file to point to the new rules locations, for example in the `directives.conf`:

```seclang
SecRuleEngine On
SecDebugLog /dev/stdout
SecDebugLogLevel 1
Include /etc/my-rules/coraza.conf-recommended
Include /etc/my-rules/crs-setup.conf.example
Include /etc/my-rules/*.conf
```

```shell
docker run \
    -v /path-to-rules:/etc/my-rules \
    -v $(pwd):/path-to-directives-file \
    ghcr.io/jcchavezs/coraza-httpbin:main \
    -directives /path-to-directives-file/directives.conf
```

### Logging

By default, coraza-httpbin will log Coraza messages to stdout. You can specify a log file by using the `-log-file` flag:

```shell
go run github.com/jcchavezs/coraza-httpbin/cmd/coraza-httpbin@latest \
  -directives ./directives.conf.example \
  -log-file coraza_log.log
```