# TraceMind 🧠

> Making sense of the chaos in distributed systems.

Distributed systems are notoriously hard to debug. When things go wrong, we're usually buried in a mountain of traces, logs, and metrics, trying to piece together a story. **TraceMind** is my attempt to automate that "story-telling" process using a hybrid approach—combining specific symbolic rules with the deductive power of LLMs.

It’s not just about "guessing" what happened; it’s about giving the AI a **Symbolic Scaffold** and a **Short-term Memory** so it can see the bigger picture.

---

## 🌩️ The Problem: "The Trace Wall"
You’ve seen it. A production incident happens, you open your observability tool, and you're met with a waterfall of 500 spans. 
*   Which one is the *actual* root cause? 
*   Is that latency spike a one-off or a trend?
*   How does the `auth` failure impact the `downstream-api`?

Usually, an SRE spends 20 minutes connecting these dots. TraceMind does it in seconds.

---

## ✨ The "Secret Sauce"

What makes TraceMind v2.0 different from a basic ChatGPT wrapper?

### 1. Hybrid Reasoning (Symbolic + Neural)
Before the LLM even sees the data, a Go-based **Symbolic Analyzer** runs a pass over the trace. It identifies bottlenecks and error origins using proven SRE heuristics. We don't just dump raw JSON into a prompt; we provide "The Facts."

### 2. Symbolic Memory
LLMs are usually stateless. TraceMind isn't. It tracks a sliding window of recent system health—error rates, slow services, and patterns. When the AI explains a trace, it knows if the system has been "shaky" for the last 10 minutes.

### 3. Real-Time SSE Streaming
Local LLMs can be slow. Instead of making you stare at a loading spinner, we stream the AI's "train of thought" live via Server-Sent Events. You watch the reasoning happen in real-time.

### 4. SRE Auditor (The Judge)
We implemented a **Research Mode** where you can compare different prompting strategies. To keep it honest, an "LLM-as-a-Judge" Auditor grades the explanations on technical accuracy and causal logic.

### 5. AI Adapter for Infrastructure Generation (NEW!)
TraceMind now includes a **Meshery AI Adapter MVP** that transforms natural language into production-ready Kubernetes infrastructure. Using the same multi-provider LLM engine, you can:

- **Generate complete infrastructure designs** from simple prompts
- **Support multiple LLM providers**: OpenAI (GPT-4), Anthropic (Claude), or local Ollama models
- **Get production-ready Kubernetes manifests** with best practices (health checks, HA, resource limits)
- **Manage AI connections** via REST API

**Example:**
```
Prompt: "Create a web app with Node.js backend, Redis cache, and PostgreSQL database"
Output: Complete Kubernetes infrastructure with Deployments, Services, StatefulSets, PVCs, ConfigMaps, and Secrets
```

See [AI_ADAPTER_README.md](AI_ADAPTER_README.md) for full documentation.

---

## 🛠️ The Build

*   **Backend**: Clean, idiomatic Go 1.25. Leveraging **LangChainGo** for orchestration.
*   **Frontend**: A sleek React + Vite dashboard designed for observability (Dark mode by default, because we're engineers).
*   **Models**: Standardized on **OTel-like** structures for future-proofing.
*   **Local-First**: Powered by **Ollama**. Your sensitive system data never leaves your machine.
*   **AI Adapter**: Multi-provider LLM support (OpenAI, Anthropic, Ollama) for infrastructure generation via REST API.

---

## � Get it Running

1.  **Ollama**: Pull the model you want to use (I recommend `mistral` or `llama3`).
    ```bash
    ollama pull mistral
    ollama serve
    ```
2.  **Backend**:
    ```bash
    export OLLAMA_MODEL=mistral
    go run main.go
    ```
    Server starts on `http://localhost:8080` with endpoints:
    - **Trace Analysis**: `/api/analyze`, `/api/evaluate`
    - **AI Connections**: `/api/connections/*`
    - **Design Generation**: `/api/design/generate`, `/api/design/generate-stream`
3.  **Frontend** (Optional):
    ```bash
    cd frontend && npm install && npm run dev
    ```

4.  **Test AI Adapter** (Optional):
    ```bash
    ./test_ai_adapter.sh
    ```

---

## 🔬 Contributing / Research

This started as a research project for the LFX Mentorship. If you're into AI-for-Ops, causal reasoning, or just want to help make system debugging less of a headache, feel free to open an issue or a PR!

The AI Adapter implementation is a **prototype MVP** of Meshery's "Natural Language to Infrastructure" capabilities. See [AI_ADAPTER_README.md](AI_ADAPTER_README.md) for architecture details and future roadmap.

---

## 📚 Documentation

- [AI Adapter API Documentation](AI_ADAPTER_README.md) - Complete guide to infrastructure generation
- [Test Script](test_ai_adapter.sh) - Automated testing for all endpoints

---

**Built with ❤️ and a lot of caffeine.**
