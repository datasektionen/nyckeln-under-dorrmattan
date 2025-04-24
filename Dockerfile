FROM golang:1.23.8-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /nyckeln

FROM alpine:3.19

COPY --from=build /nyckeln /nyckeln
COPY config.yaml ./

CMD ["/nyckeln"]

LABEL org.opencontainers.image.source="https://github.com/datasektionen/nyckeln-under-dorrmattan" \
      org.opencontainers.image.licenses="MIT"
