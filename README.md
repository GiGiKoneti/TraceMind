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

---

## 🛠️ The Build

*   **Backend**: Clean, idiomatic Go 1.25. Leveraging **LangChainGo** for orchestration.
*   **Frontend**: A sleek React + Vite dashboard designed for observability (Dark mode by default, because we're engineers).
*   **Models**: Standardized on **OTel-like** structures for future-proofing.
*   **Local-First**: Powered by **Ollama**. Your sensitive system data never leaves your machine.

---

## � Get it Running

1.  **Ollama**: Pull the model you want to use (I recommend `mistral` or `llama3`).
    ```bash
    ollama pull mistral
    ```
2.  **Backend**:
    ```bash
    export OLLAMA_MODEL=mistral
    go run main.go
    ```
3.  **Frontend**:
    ```bash
    cd frontend && npm install && npm run dev
    ```

---

## 🔬 Contributing / Research
This started as a research project for the LFX Mentorship. If you're into AI-for-Ops, causal reasoning, or just want to help make system debugging less of a headache, feel free to open an issue or a PR!

**Built with ❤️ and a lot of caffeine.**
