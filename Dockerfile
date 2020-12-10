FROM golang:1.14.6-alpine as builder

WORKDIR /build/

RUN apk update && \
    apk upgrade && \
    apk add ca-certificates git openssh-client

COPY . .

RUN go build -o bin/server cmd/api/main.go


FROM alpine:3.7

WORKDIR /opt/test-document/
RUN apk --no-cache add bash ca-certificates tzdata
COPY --from=builder /build/bin/server bin/server


CMD bin/server