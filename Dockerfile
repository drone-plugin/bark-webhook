FROM alpine:latest

ADD  main /bin/
ENTRYPOINT /bin/main