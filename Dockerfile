FROM golang:1.23

ENV GO111MODULE=on
ENV PORT=8000
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./api

EXPOSE 8000
ENTRYPOINT ["./api"] 