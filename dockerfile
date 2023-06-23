# Start from the golang image
FROM golang:alpine as compiler

# Set the working directory for golang image
WORKDIR /app

# Copy go.mod and go.sum files to the /app directory
COPY go.mod .
COPY go.sum .

# Install dependencies for our project using Go modules
RUN go mod download

# Copy all project files inside the container to /app
COPY . .

# Build our application
RUN go build -o golang-project

#specify the base image for the Docker container as Alpine Linux
FROM alpine

# Set the working directory for our application
WORKDIR /small

# Copy all project files from golang working directory to Alpine Linux Docker container 
COPY --from=compiler ./app/golang-project ./binary

# Specify the initial command that should be executed
ENTRYPOINT [ "./binary" ]

