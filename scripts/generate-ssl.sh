#!/bin/bash

# Generate self-signed SSL certificate for development
# For production, use Let's Encrypt or a proper CA

mkdir -p nginx/ssl

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout nginx/ssl/key.pem \
    -out nginx/ssl/cert.pem \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

echo "SSL certificates generated in nginx/ssl/"
echo "For production, replace these with proper certificates from a CA"
