FROM golang:1.21.5 as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go mod and sum files to leverage Docker cache.
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

COPY . .

# Navigate to the cmd directory.
WORKDIR /app/cmd

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio-backend

# Use minimal base image.

FROM golang:1.21.5-alpine

# Set the RUNNING_IN_DOCKER environment variable
ENV RUNNING_IN_DOCKER=true

COPY --from=builder /app/cmd/portfolio-backend /app/portfolio-backend

# Copy necessary resources

COPY --from=builder /app/pkg /app/pkg

# Set the working directory
WORKDIR /app

# Expose the port on which the app runs.
EXPOSE 8080

CMD [ "./portfolio-backend" ]