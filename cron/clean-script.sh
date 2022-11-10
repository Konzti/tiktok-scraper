#!/bin/bash
path="/tiktok/"
echo "trying to delete files..."
sleep 1
find $path -mmin +5 -type f -print -delete