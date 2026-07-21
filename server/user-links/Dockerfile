# syntax=docker/dockerfile:1

FROM golang:1.26.5
LABEL authors="github.com/7om3k"

WORKDIR /app

# Docker's COPY command is quite sensitive in its interpretation of the trailing slash
COPY ./applib/ ./applib/

COPY ./user-links/go.mod ./user-links/
COPY ./user-links/*.go ./user-links/

RUN go work init ./applib ./user-links

WORKDIR /app/user-links
RUN go mod download

# Binary is located in the root of the filesystem of the image
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-user-links

CMD ["/app-user-links"]
