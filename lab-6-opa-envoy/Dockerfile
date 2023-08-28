FROM golang:1.20.7-alpine3.18 as builder

WORKDIR build

COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o employee-records

# Create a new release build stage
FROM scratch

# Set the working directory to the root directory path
WORKDIR /

# Copy over the binary built from the previous stage
COPY --from=builder /go/build/employee-records /employee-records

EXPOSE 8080

ENTRYPOINT ["/employee-records"]
