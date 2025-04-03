FROM alpine:latest

WORKDIR /app

# Copy backend and frontend binaries into the container
COPY backend/image /app/backend
COPY frontend/image /app/frontend

# Expose backend port
EXPOSE 8080

# Start both processes
CMD ["/bin/sh", "-c", "/app/backend & /app/frontend"]