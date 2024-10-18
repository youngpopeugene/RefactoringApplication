FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/cmd/app
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]