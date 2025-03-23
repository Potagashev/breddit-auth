FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# there is the code to run in dev
CMD ["go", "run", "main.go"]

# there is the code to run in prod (maybe...)
# COPY . .
# RUN go build -o main .

# FROM alpine
# WORKDIR /root/
# COPY --from=builder /app/main .
# ADD . /app

# CMD ["./main"]