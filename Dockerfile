FROM golang:latest as builder
LABEL maintainer="Thomas Gorham <tgiv014@gmail.com>"

# Install npm
RUN apt-get update && apt-get install -y npm

# Install task and templ
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app
COPY go.mod go.sum package.json package-lock.json ./
RUN go mod download
RUN npm install

COPY . .

RUN task build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/to /app/to

EXPOSE 8080

CMD ["/app/to"]