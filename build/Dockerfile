# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-alpine AS build

ARG BUILD_OS_TARGET
ARG BUILD_ARCH_TARGET

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY config ./config
COPY internal ./internal

RUN GOOS=$BUILD_OS_TARGET GOARCH=$BUILD_ARCH_TARGET CGO_ENABLED=0 go build -o /transmission-web-scrapper ./cmd/transmission-web-scrapper/main.go

CMD [ "/transmission-web-scrapper" ]

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /transmission-web-scrapper /transmission-web-scrapper
COPY --from=build /app/config/.env.local /app/config/.env.prod /config/

ENTRYPOINT ["/transmission-web-scrapper"]