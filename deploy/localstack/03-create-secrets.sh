#!/bin/bash
echo "Initializing LocalStack Secrets Manager..."

awslocal secretsmanager create-secret --name donald/jwt-secret --secret-string "super-secret-jwt-key"
awslocal secretsmanager create-secret --name donald/db-password --secret-string "local-db-password"

echo "Secrets created."
