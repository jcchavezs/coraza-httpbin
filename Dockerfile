FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY ./build/coraza-httpbin-${TARGETOS}-${TARGETARCH} /usr/bin/coraza-httpbin
COPY ./directives.conf.example /etc/coraza-httpbin/directives.conf.example

EXPOSE 8080

CMD ["-directives", "/etc/coraza-httpbin/directives.conf.example"]

ENTRYPOINT [ "/usr/bin/coraza-httpbin", "-port", "8080" ]