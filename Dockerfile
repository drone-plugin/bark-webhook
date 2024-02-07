FROM alpine:latest

ADD  main /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/main