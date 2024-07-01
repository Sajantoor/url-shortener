FROM golang:1.22-bullseye as builder 

WORKDIR /app

ADD . /app

RUN go build -o bin/app services/creation/main.go

FROM debian:bullseye-slim

COPY --from=builder /app/bin/app . 

EXPOSE 8080 

CMD ["./app"]
