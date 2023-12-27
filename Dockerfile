FROM golang:1.21-alpine as builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY ./go.mod ./
RUN go mod download
COPY . .
RUN templ generate -path="./view"
RUN CGO_ENABLED=0 go build -o server cmd/shah/main.go
EXPOSE 8080
CMD ["/app/server"]