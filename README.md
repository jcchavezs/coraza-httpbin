# coraza-httpbin

`httpbin` server with coraza as reverse proxy to play with attacks.

## Getting started

```shell
go run github.com/jcchavezs/coraza-httpbin/cmd/coraza-httpbin@latest
```

And do a call in another terminal:

```shell
# attemp XSS injection
curl -i "localhost:8080?arg=<script>a</script>"
```
