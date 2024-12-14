FROM golang:1.19.1-alpine AS builder
ARG VERSION=dev

ENV APP_HOME /go/src/kingwood-api
WORKDIR "$APP_HOME"

COPY . .
COPY ./.env .

# RUN go mod download
# RUN go mod verify

RUN CGO_ENABLED=0 go build -o main -ldflags=-X=main.version=${VERSION} cmd/main.go

FROM alpine:3.14
LABEL org.opencontainers.image.source=https://github.com/mikalai2006/kingwood-api
LABEL org.opencontainers.image.description="Template REST API"
LABEL org.opencontainers.image.licenses=MIT

ENV APP_HOME /go/src/kingwood-api
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY configs/ configs/
COPY --from=builder "$APP_HOME"/.env /go/bin/.env
COPY --from=builder "$APP_HOME" /go/bin

EXPOSE 8000
ENV PATH="/go/bin:${PATH}"
CMD ["main"]