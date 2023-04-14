# coraza-httpbin

`httpbin` server with coraza as reverse proxy to play with attacks.

## Getting started

### Running in Go

First you need to run the server using a directives file. Notice `@owasp_crs` folder
is already included.

```shell
go run github.com/jcchavezs/coraza-httpbin/cmd/coraza-httpbin@latest -directives ./directives.conf.example
```

and do a HTTP call in another terminal:

```shell
# attempt XSS injection
curl -i "localhost:8080?arg=<script>a</script>"
```

You can also point to rules in your filesystem by using the absolute notation in the directives file:

```seclang
SecRuleEngine On
SecDebugLog /dev/stdout
SecDebugLogLevel 1
Include /path-to-my-rules/crs-setup.conf.example
Include /path-to-my-rules/*.conf
```

### Running in Docker

```shell
docker run -v /path-to-rules:/etc/my-rules ghcr.io/jcchavezs/coraza-httpbin:main
```

and change your directives file to point to the new rules locations, for example in the `directives.conf.example`:

```seclang
SecRuleEngine On
SecDebugLog /dev/stdout
SecDebugLogLevel 1
Include /etc/my-rules/crs-setup.conf.example
Include /etc/my-rules/*.conf
```

```shell
docker run -v /path-to-rules:/etc/my-rules ghcr.io/jcchavezs/coraza-httpbin:main -directives /path-to-directives-file.conf
```
