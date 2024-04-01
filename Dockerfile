FROM golang:1.22-alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./
COPY login login
COPY pls pls

RUN CGO_ENABLED=0 GOOS=linux go build -o /nyckeln

FROM alpine:3.19

COPY --from=build /nyckeln /nyckeln

CMD ["/nyckeln"]

LABEL org.opencontainers.image.source="https://github.com/datasektionen/nyckeln-under-dorrmattan" \
      org.opencontainers.image.licenses="MIT"
