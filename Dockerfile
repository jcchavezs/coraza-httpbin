FROM scratch

COPY build/coraza-httpbin-linux /usr/bin/coraza-httpbin

EXPOSE 8080

CMD [ "-port", "8080"]

ENTRYPOINT [ "/usr/bin/coraza-httpbin" ]