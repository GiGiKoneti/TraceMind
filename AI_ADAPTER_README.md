# Meshery AI Adapter MVP - API Documentation

## Overview

This prototype implements a simplified version of Meshery's AI Adapter for natural language to infrastructure design generation. It supports multiple LLM providers (OpenAI, Anthropic, Ollama) and provides REST APIs for connection management and design generation.

## Quick Start

### Prerequisites

1. **Go 1.25+** installed
2. **Ollama** running locally (for local LLM support)
   ```bash
   ollama serve
   ollama pull mistral
   ```
3. **(Optional)** OpenAI or Anthropic API keys for cloud LLM providers

### Running the Server

```bash
# Start the server (defaults to Ollama with mistral model)
export OLLAMA_MODEL=mistral
go run main.go

# Or build and run
go build -o tracemind-ai .
./tracemind-ai
```

Server starts on `http://localhost:8080` by default.

---

## API Endpoints

### Connection Management

#### 1. Create AI Connection

**POST** `/api/connections/create`

Create a new AI provider connection.

**Request Body:**
```json
{
  "name": "Connection Name",
  "provider": "openai|anthropic|ollama",
  "config": {
    "model": "model-name",
    "endpoint": "optional-endpoint",
    "api_endpoint": "optional-api-endpoint"
  },
  "credentials": {
    "api_key": "your-api-key"
  }
}
```

**Examples:**

**Ollama (Local):**
```bash
curl -X POST http://localhost:8080/api/connections/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Local Mistral",
    "provider": "ollama",
    "config": {
      "endpoint": "http://localhost:11434",
      "model": "mistral"
    },
    "credentials": {}
  }'
```

**OpenAI:**
```bash
curl -X POST http://localhost:8080/api/connections/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "OpenAI GPT-4",
    "provider": "openai",
    "config": {
      "api_endpoint": "https://api.openai.com/v1",
      "model": "gpt-4"
    },
    "credentials": {
      "api_key": "sk-..."
    }
  }'
```

**Anthropic:**
```bash
curl -X POST http://localhost:8080/api/connections/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Claude Opus",
    "provider": "anthropic",
    "config": {
      "model": "claude-3-opus-20240229"
    },
    "credentials": {
      "api_key": "sk-ant-..."
    }
  }'
```

**Response:**
```json
{
  "id": "uuid",
  "name": "Connection Name",
  "provider": "openai",
  "config": {...},
  "status": "untested",
  "created_at": "2026-02-11T02:00:00Z"
}
```

---

#### 2. List Connections

**GET** `/api/connections`

List all AI connections.

```bash
curl http://localhost:8080/api/connections
```

**Response:**
```json
[
  {
    "id": "uuid-1",
    "name": "Local Mistral",
    "provider": "ollama",
    "status": "connected",
    "created_at": "2026-02-11T02:00:00Z"
  },
  {
    "id": "uuid-2",
    "name": "OpenAI GPT-4",
    "provider": "openai",
    "status": "connected",
    "created_at": "2026-02-11T02:05:00Z"
  }
]
```

---

#### 3. Test Connection

**POST** `/api/connections/test`

Test if a connection is working.

```bash
curl -X POST http://localhost:8080/api/connections/test \
  -H "Content-Type: application/json" \
  -d '{
    "connection_id": "uuid"
  }'
```

**Response (Success):**
```json
{
  "status": "connected",
  "message": "Connection successful",
  "response": "OK"
}
```

**Response (Error):**
```json
{
  "status": "error",
  "message": "connection refused"
}
```

---

#### 4. Delete Connection

**POST** `/api/connections/delete`

Delete an AI connection.

```bash
curl -X POST http://localhost:8080/api/connections/delete \
  -H "Content-Type: application/json" \
  -d '{
    "connection_id": "uuid"
  }'
```

---

### Design Generation

#### 5. Generate Infrastructure Design

**POST** `/api/design/generate`

Generate infrastructure design from natural language prompt.

```bash
curl -X POST http://localhost:8080/api/design/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Deploy a highly available Kubernetes cluster with Prometheus monitoring and Grafana dashboards",
    "connection_id": "uuid"
  }'
```

**Response:**
```json
{
  "name": "k8s-monitoring-stack",
  "description": "Highly available Kubernetes cluster with monitoring",
  "version": "1.0.0",
  "components": [
    {
      "id": "uuid",
      "name": "prometheus-deployment",
      "type": "Kubernetes",
      "apiVersion": "apps/v1",
      "kind": "Deployment",
      "spec": {
        "replicas": 3,
        "selector": {...},
        "template": {...}
      },
      "metadata": {
        "namespace": "monitoring",
        "labels": {...}
      }
    },
    {
      "id": "uuid",
      "name": "grafana-service",
      "type": "Kubernetes",
      "apiVersion": "v1",
      "kind": "Service",
      "spec": {...},
      "metadata": {...}
    }
  ],
  "metadata": {
    "generated_by": "TraceMind AI Adapter",
    "prompt": "Deploy a highly available...",
    "provider": "openai",
    "model": "gpt-4",
    "generated_at": "2026-02-11T02:10:00Z"
  }
}
```

---

#### 6. Generate Design (Streaming)

**POST** `/api/design/generate-stream`

Generate infrastructure design with real-time streaming.

```bash
curl -X POST http://localhost:8080/api/design/generate-stream \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Create a microservices architecture with API gateway, service mesh, and distributed tracing",
    "connection_id": "uuid"
  }'
```

**Response (Server-Sent Events):**
```
event: token
data: {

event: token
data: "name"

event: token
data: :

event: design
data: {"name":"microservices-stack",...}

event: done
data: [DONE]
```

---

## Example Prompts

Here are some example prompts to try:

1. **Simple Web App:**
   ```
   Create a simple web application with a Node.js backend, Redis cache, and PostgreSQL database. Include health checks.
   ```

2. **Microservices:**
   ```
   Deploy a microservices architecture with 3 services (auth, api, worker), Nginx ingress, and RabbitMQ message queue.
   ```

3. **Data Pipeline:**
   ```
   Build a data processing pipeline with Kafka for streaming, Spark for processing, and Elasticsearch for storage.
   ```

4. **Monitoring Stack:**
   ```
   Set up complete observability with Prometheus, Grafana, Loki for logs, and Jaeger for distributed tracing.
   ```

5. **ML Platform:**
   ```
   Create a machine learning platform with JupyterHub, MLflow for experiment tracking, and MinIO for model storage.
   ```

---

## Testing

### Manual Testing

Use the provided test script:

```bash
./test_ai_adapter.sh
```

### Testing with Ollama (Local)

```bash
# 1. Start Ollama
ollama serve

# 2. Pull a model
ollama pull mistral

# 3. Start TraceMind
go run main.go

# 4. Create Ollama connection
curl -X POST http://localhost:8080/api/connections/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Local Mistral",
    "provider": "ollama",
    "config": {
      "endpoint": "http://localhost:11434",
      "model": "mistral"
    },
    "credentials": {}
  }'

# 5. Generate design (replace CONNECTION_ID)
curl -X POST http://localhost:8080/api/design/generate \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Create a simple Redis cache with persistent storage",
    "connection_id": "CONNECTION_ID"
  }' | jq '.'
```

---

## Architecture

```
TraceMind AI Adapter
├── internal/
│   ├── models/
│   │   ├── connections.go    # AI connection models
│   │   ├── meshery.go        # Infrastructure design models
│   │   └── trace.go          # Existing trace models
│   ├── llm/
│   │   ├── engine.go         # LLM engine (backward compatible)
│   │   ├── provider.go       # Provider interface & factory
│   │   ├── openai_provider.go
│   │   ├── anthropic_provider.go
│   │   └── ollama_provider.go
│   └── handlers/
│       ├── api.go            # Existing trace handlers
│       ├── connections.go    # Connection management
│       └── ai_design.go      # Design generation
└── main.go                   # Server with all routes
```

---

## Differences from Full Meshery Implementation

This is a **prototype MVP** with the following simplifications:

1. **No gRPC Adapter Pattern** - Uses REST API only
2. **No GraphQL** - REST endpoints only
3. **In-Memory Storage** - No database persistence
4. **Simplified Design Models** - Not full Meshery Models registry
5. **No Frontend UI** - API only (can be tested with curl/Postman)

---

## Next Steps

To evolve this into a full Meshery integration:

1. Add database persistence (PostgreSQL/SQLite)
2. Implement gRPC adapter pattern
3. Add GraphQL API layer
4. Integrate with Meshery Models registry
5. Add Kubernetes manifest validation
6. Implement rate limiting and cost tracking
7. Build frontend UI components
8. Add comprehensive test suite

---

## Troubleshooting

### "Connection refused" for Ollama

Make sure Ollama is running:
```bash
ollama serve
```

### "Model not found"

Pull the model first:
```bash
ollama pull mistral
```

### OpenAI/Anthropic API errors

Verify your API key is valid and has sufficient credits.

---

## License

MIT License (same as TraceMind)
