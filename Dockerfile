FROM golang:1.20

WORKDIR /app
COPY main.go go.mod ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /login

CMD ["/login"]
