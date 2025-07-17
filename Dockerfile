FROM golang:1.24.4-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build CGO_ENABLED=0 GOOS=linux -o /main

FROM alpine:latest AS run

WORKDIR /

COPY --from=build /main /main
COPY --from=build /app/migrations /migrations

ENTRYPOINT [ "/main" ]
