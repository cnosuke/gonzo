FROM alpine:3.8

RUN apk --update add --no-cache ca-certificates

WORKDIR /app
COPY ./dist/imagine ./

RUN mkdir -p ./static/
COPY ./_front/dist ./static/

CMD ["./imagine", "-v"]
