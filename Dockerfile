FROM golang:latest as builder
LABEL maintainer="Adam Jordan <adamyordan@gmail.com>"
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o laplace .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/laplace .
COPY files files
EXPOSE 443
CMD ["./laplace"]
