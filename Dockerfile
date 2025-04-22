# Builder stage
FROM golang:1.23.8 AS builder

# Install needed tools (might not be strictly necessary for just building pixterm, but kept from original)
RUN apt-get update && apt-get install -y unzip curl

# Set up working dir
WORKDIR /build

# Download and extract the release zip
RUN curl -L https://github.com/eliukblau/pixterm/archive/refs/tags/v1.3.2.zip -o pixterm.zip \
    && unzip pixterm.zip

# Build the binary
WORKDIR /build/pixterm-1.3.2/cmd/pixterm
RUN go build -o pixterm .

# --- End of Builder Stage ---2

ENTRYPOINT ["/build/pixterm-1.3.2/cmd/pixterm/pixterm"]


# Runtime stage
# Using a minimal base image like alpine
# FROM alpine:latest

# # Create a directory for the application

# # Copy the compiled binary from the builder stage into the runtime image

# # Set the working directory to /app
# WORKDIR /app
# COPY --from=builder /build/pixterm-1.3.2/cmd/pixterm/pixterm /app/pixterm
# # Set the entrypoint to the pixterm binary.
# # This ensures that when the container runs, it executes ./pixterm
# # Any command-line arguments provided to 'docker run' will be passed as arguments to ./pixterm
# ENTRYPOINT ["/app/pixterm"]

# You could optionally set a default command/arguments if no args are provided to docker run
# For example, to show help by default:
# CMD ["--help"]