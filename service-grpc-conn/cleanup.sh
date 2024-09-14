#!/bin/bash

if [ -f "./commit-log" ]; then
    echo "Removing existing commit-log executable..."
    rm ./commit-log
fi

echo "Removing hidden files..."
find . -name "._*" -exec rm -v {} \;

if [ -d "./log_data" ]; then
    echo "Removing log_data directory..."
    rm -rf ./log_data
fi
