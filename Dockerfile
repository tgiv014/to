FROM golang:latest as builder
LABEL maintainer="Thomas Gorham <tgiv014@gmail.com>"


# Install task and templ
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN task build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]