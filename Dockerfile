FROM golang:1.21-alpine as builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
COPY ./go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/shah/main.go
EXPOSE 8080
CMD ["/app/server"]