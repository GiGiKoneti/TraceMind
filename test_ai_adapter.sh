#!/bin/bash

# TraceMind AI Adapter - Test Script
# This script demonstrates the AI connection and design generation capabilities

BASE_URL="http://localhost:8080"

echo "=== TraceMind AI Adapter Test Suite ==="
echo ""

# Test 1: Create Ollama Connection
echo "Test 1: Creating Ollama Connection..."
curl -X POST "${BASE_URL}/api/connections/create" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Local Ollama Mistral",
    "provider": "ollama",
    "config": {
      "endpoint": "http://localhost:11434",
      "model": "mistral"
    },
    "credentials": {}
  }' | jq '.'

echo -e "\n"

# Test 2: List All Connections
echo "Test 2: Listing All Connections..."
curl -X GET "${BASE_URL}/api/connections" | jq '.'

echo -e "\n"

# Test 3: Create OpenAI Connection (requires API key)
echo "Test 3: Creating OpenAI Connection (set OPENAI_API_KEY env var)..."
if [ -n "$OPENAI_API_KEY" ]; then
  curl -X POST "${BASE_URL}/api/connections/create" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"OpenAI GPT-4\",
      \"provider\": \"openai\",
      \"config\": {
        \"api_endpoint\": \"https://api.openai.com/v1\",
        \"model\": \"gpt-4\"
      },
      \"credentials\": {
        \"api_key\": \"$OPENAI_API_KEY\"
      }
    }" | jq '.'
else
  echo "Skipping - OPENAI_API_KEY not set"
fi

echo -e "\n"

# Test 4: Test Connection (get connection ID from list)
echo "Test 4: Testing Connection..."
CONNECTION_ID=$(curl -s "${BASE_URL}/api/connections" | jq -r '.[0].id')
if [ -n "$CONNECTION_ID" ] && [ "$CONNECTION_ID" != "null" ]; then
  curl -X POST "${BASE_URL}/api/connections/test" \
    -H "Content-Type: application/json" \
    -d "{
      \"connection_id\": \"$CONNECTION_ID\"
    }" | jq '.'
else
  echo "No connections available to test"
fi

echo -e "\n"

# Test 5: Generate Infrastructure Design
echo "Test 5: Generating Infrastructure Design..."
if [ -n "$CONNECTION_ID" ] && [ "$CONNECTION_ID" != "null" ]; then
  curl -X POST "${BASE_URL}/api/design/generate" \
    -H "Content-Type: application/json" \
    -d "{
      \"prompt\": \"Create a simple web application with a Node.js backend, Redis cache, and PostgreSQL database. Include health checks and monitoring.\",
      \"connection_id\": \"$CONNECTION_ID\"
    }" | jq '.'
else
  echo "No connections available for design generation"
fi

echo -e "\n"
echo "=== Test Suite Complete ==="
