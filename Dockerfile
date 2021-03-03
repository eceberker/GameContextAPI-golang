# FROM golang:alpine as builder
# COPY go.mod go.sum /go/src/github.com/eceberker/gamecontextdb/
# WORKDIR /go/src/github.com/eceberker/gamecontextdb
# RUN go mod download
# COPY . /go/src/github.com/eceberker/gamecontextdb
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/gamecontextdb github.com/eceberker/gamecontextdb
# Start from golang base image
FROM golang:alpine as builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
# Set the current working directory inside the container 
WORKDIR /app
# Copy go mod and sum files 
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 
# Copy the source from the current directory to the working Directory inside the container 
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# Expose port 8080 to the outside world
EXPOSE 8080
#Command to run the executable
CMD ["./main"]

# FROM alpine
# RUN apk add --no-cache ca-certificates && update-ca-certificates
# COPY --from=builder /go/src/github.com/eceberker/gamecontextdb/build/gamecontextdb /usr/bin/gamecontextdb
# EXPOSE 8080 8080
# ENTRYPOINT ["/usr/bin/gamecontextdb"]