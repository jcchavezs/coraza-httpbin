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

By default, coraza-httpbin will log Coraza messages to stdout. You can specify a log file by using the `-matched-rules-log` flag:

```shell
go run github.com/jcchavezs/coraza-httpbin/cmd/coraza-httpbin@latest \
  -directives ./directives.conf.example \
  -matched-rules-log coraza_log.log
```

and get logs:

```
[critical] [client "127.0.0.1"] Coraza: Warning. XSS Attack Detected via libinjection [file "../coraza-httpbin/@owasp_crs/REQUEST-941-APPLICATION-ATTACK-XSS.conf"] [line "4434"] [id "941100"] [rev ""] [msg "XSS Attack Detected via libinjection"] [data "Matched Data: XSS data found within ARGS:arg: <script>a</script>"] [severity "critical"] [ver "OWASP_CRS/4.0.0-rc1"] [maturity "0"] [accuracy "0"] [tag "application-multi"] [tag "language-multi"] [tag "platform-multi"] [tag "attack-xss"] [tag "paranoia-level/1"] [tag "OWASP_CRS"] [tag "capec/1000/152/242"] [hostname ""] [uri "/?arg=<script>a</script>"] [unique_id "GRdXJCfkPuFbQwyLXlF"]
```
