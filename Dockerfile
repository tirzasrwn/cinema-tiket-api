FROM golang:1.25.0-alpine3.21 AS builder
LABEL stage=builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN sed -i 's/localhost/cinema-ticket-api-postgres/g' ./.env
RUN swag init -g ./main.go --output ./docs --parseDependency
USER 0:0
RUN go build -o ./backend .
RUN chmod -R 777 ./backend

FROM alpine:3.21 AS production
WORKDIR /app
COPY --from=builder /app ./
CMD ["./backend"]
EXPOSE 4000 

