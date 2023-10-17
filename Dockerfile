FROM golang:1.21-bookworm as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

# Build the Go web service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pastebin ./cmd

# Use the most lightweight image for the final runtime
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/pastebin .

EXPOSE 50998

CMD ["./pastebin", "--port", "50998"]
