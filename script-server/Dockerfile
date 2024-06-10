FROM alpine:latest

# Install socat
RUN apk add --no-cache socat

# Create a directory for the script
RUN mkdir -p /usr/local/scripts

# Copy the script into the image
COPY script.sh /usr/local/scripts/script.sh

# Expose port 8080 for the HTTP server
EXPOSE 8080

# Set the working directory
WORKDIR /usr/local/scripts

# Create a simple HTTP response script that handles signals
RUN echo -e "#!/bin/sh\n\
trap \"exit 0\" SIGTERM SIGINT\n\
while :; do\n\
  { \n\
    echo -e 'HTTP/1.1 200 OK\\r\\nContent-Type: text/plain\\r\\n\\r\\n'; \n\
    cat /usr/local/scripts/script.sh; \n\
  } | socat - TCP-LISTEN:8080,forever,reuseaddr\n\
done" > serve.sh && chmod +x serve.sh

# Start serving the script
CMD ["./serve.sh"]
