FROM scratch

COPY build/coraza-httpbin-linux /usr/bin/coraza-httpbin

COPY ./directives.conf.example /etc/coraza-httpbin/directives.conf.example

EXPOSE 8080

CMD ["-directives", "/etc/coraza-httpbin/directives.conf.example"]

ENTRYPOINT [ "/usr/bin/coraza-httpbin", "-port", "8080" ]