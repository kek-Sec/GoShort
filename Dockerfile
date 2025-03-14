# Stage 1: Build the Go application
FROM golang:1.22.4 AS builder

ARG VERSION=dev
ENV VERSION=${VERSION}

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=${VERSION}" -o goshort ./cmd/server

# Stage 2: Build the SvelteKit web application
FROM node:18 AS web-builder

# Set working directory
WORKDIR /web

# Copy the SvelteKit project files
COPY web/package.json web/package-lock.json ./
COPY web/ ./

# Install dependencies and build the web app
RUN npm install
RUN npm run build

# Stage 3: Production
FROM alpine:latest AS production

# Install Nginx and Supervisord
RUN apk add --no-cache nginx supervisor

# Set up working directories
WORKDIR /app

# Copy the Go application binary
COPY --from=builder /app/goshort /app/goshort

# Copy the SvelteKit build output to Nginx HTML folder
COPY --from=web-builder /web/build /usr/share/nginx/html

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy supervisord configuration
COPY supervisord.conf /etc/supervisord.conf

# Set environment variables with default values that can be overridden
ENV BRAND_TITLE="GoShort - URL Shortener" \
    BRAND_DESCRIPTION="GoShort is a powerful and user-friendly URL shortener. Simplify, manage, and track your links with ease." \
    BRAND_KEYWORDS="URL shortener, GoShort, link management, shorten URLs, track links" \
    BRAND_AUTHOR="GoShort Team" \
    BRAND_THEME_COLOR="#4caf50" \
    BRAND_LOGO_TEXT="GoShort" \
    BRAND_PRIMARY_COLOR="#3b82f6" \
    BRAND_SECONDARY_COLOR="#10b981" \
    BRAND_HEADER_TITLE="GoShort - URL Shortener" \
    BRAND_FOOTER_TEXT="View the project on" \
    BRAND_FOOTER_LINK="https://github.com/kek-Sec/GoShort"

# Expose ports
EXPOSE 80 8080

# Command to run both Nginx and the Go application
CMD ["supervisord", "-c", "/etc/supervisord.conf"]
