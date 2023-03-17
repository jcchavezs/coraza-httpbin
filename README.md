# coraza-httpbin

`httpbin` server with coraza as reverse proxy to play with attacks.

## Getting started

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
