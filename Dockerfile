FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o url-api ./cmd/main.go

EXPOSE 5050

CMD [ "./url-api" ]