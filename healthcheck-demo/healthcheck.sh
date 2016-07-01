#!/bin/sh
echo "Executing check..."
if [ ! -f "/opt/healthcheck" ]; then
    echo "File not found!"
    exit 1
fi
