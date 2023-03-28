FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY /build/coraza-httpbin-${TARGETOS}-${TARGETARCH} /usr/bin/coraza-httpbin

EXPOSE 8080

CMD [ "-port", "8080"]

ENTRYPOINT [ "/usr/bin/coraza-httpbin" ]