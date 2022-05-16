FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /astral-api

EXPOSE 3000

CMD [ "/astral-api" ]