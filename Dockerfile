FROM golang:1.21.6-alpine as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go mod download
RUN go mod verify
RUN go build -v cmd/paste/main.go

FROM golang:1.21.6-alpine
COPY --from=build /usr/src/app/main .
COPY static ./static

CMD ["./main"]
