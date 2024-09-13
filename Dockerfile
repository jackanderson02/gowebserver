FROM golang:1.23.0-alpine3.19

RUN mkdir /app

ADD . /app

WORKDIR /app

EXPOSE 8080

RUN go build -o main .

CMD ["/app/main"]