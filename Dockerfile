FROM alpine:3.8

RUN apk --update add --no-cache ca-certificates

WORKDIR /app
COPY ./dist/gonzo ./

CMD ["./gonzo", "-v"]
