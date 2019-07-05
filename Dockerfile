FROM alpine:3.8

RUN apk --update add --no-cache ca-certificates

WORKDIR /app
COPY ./dist/gonzo ./

RUN mkdir -p /etc
COPY mime.types /etc/mime.types

CMD ["./gonzo", "-v"]
