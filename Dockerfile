FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .env .
COPY . .

RUN go build -o /ozinshe-go 

EXPOSE 8080

ENTRYPOINT ["/ozinshe-go"]