# syntax=docker/dockerfile:1

FROM golang:1.27.1
LABEL authors="github.com/7om3k"

WORKDIR /app

# Docker's COPY command is quite sensitive in its interpretation of the trailing slash
COPY ./applib ./applib/

COPY ./web-ui/pages ./web-ui/pages/
COPY ./web-ui/go.mod ./web-ui/
COPY ./web-ui/*.go ./web-ui/

RUN go work init ./applib ./web-ui

WORKDIR /app/web-ui

# Binary is located in the root of the filesystem of the image
RUN CGO_ENABLED=0 GOOS=linux go build -o /web-ui

CMD ["/web-ui"]
