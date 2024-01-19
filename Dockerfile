# Use an official Go runtime as a parent image
FROM golang:1.18


# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Switch back to the root directory


# Copy the necessary files for the frontend
COPY ./frontend /frontend

# Expose the port on which your Go application runs
EXPOSE 8001
EXPOSE 3306

# Command to run the executable
CMD ["./main"]
