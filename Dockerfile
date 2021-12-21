FROM golang:1.17.1-alpine as builder
WORKDIR /build

COPY . /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o libraryBin .
RUN go mod download

# generate clean, final image for end users
FROM alpine:3.11.3
WORKDIR /library
COPY --from=builder /build/ /library/

# executable
ENTRYPOINT [ "/library/libraryBin" ]