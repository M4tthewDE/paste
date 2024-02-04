FROM golang:1.21.6-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY cmd ./cmd
COPY internal ./internal

RUN go build -v cmd/paste/main.go

CMD ["./main"]
