#!/bin/bash
echo "Initializing LocalStack S3..."

awslocal s3 mb s3://donald-vibe
awslocal s3api put-bucket-versioning --bucket donald-vibe --versioning-configuration Status=Enabled
awslocal s3api put-bucket-acl --bucket donald-vibe --acl public-read

echo "S3 Bucket 'donald-vibe' created."
