#!/bin/bash
echo "Initializing LocalStack SQS..."

awslocal sqs create-queue --queue-name donald-orders
awslocal sqs create-queue --queue-name donald-emails
awslocal sqs create-queue --queue-name donald-notifications

echo "SQS Queues created."
