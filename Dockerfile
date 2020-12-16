FROM golang:1.15.2 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zrule ./cmd/zrule

FROM alpine:latest AS release
WORKDIR /app

RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /app/zrule .

LABEL maintainer="David Douglas <david@onetwentyseven.dev>"