#!/bin/bash

# Check if two IP addresses are provided as arguments
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <vm1-ip> <vm2-ip>"
  exit 1
fi

# Assign provided IP addresses to variables
VM1_IP="$1"
VM2_IP="$2"

# Specify the file path
FILE_PATH="ansible/hosts"

# Check if the file exists
if [ ! -f "$FILE_PATH" ]; then
  echo "Error: File '$FILE_PATH' not found."
  exit 1
fi

# Replace placeholders with actual IP addresses in the file
sed -i "s/<vm1-ip>/$VM1_IP/g" "$FILE_PATH"
sed -i "s/<vm2-ip>/$VM2_IP/g" "$FILE_PATH"

echo "IP addresses replaced successfully in '$FILE_PATH'."
