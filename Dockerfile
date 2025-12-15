FROM golang:1.23-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /licenseplate main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /licenseplate /usr/local/bin/licenseplate
WORKDIR /app
ENV PORT=9002
EXPOSE 9002
CMD ["/usr/local/bin/licenseplate"]
