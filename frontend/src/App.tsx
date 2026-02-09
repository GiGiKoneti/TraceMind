import { useState, useRef } from 'react'

interface Span {
  span_id: string
  name: string
  start_time: string
  end_time: string
  status: { code: string; message?: string }
}

interface Trace {
  trace_id: string
  spans: Span[]
}

interface SymbolicFact {
  type: string
  service: string
  description: string
  severity: 'info' | 'warning' | 'critical'
}

interface SystemHealth {
  recent_error_rate: number
  slowest_services: string[]
  last_update: string
}

interface AnalysisState {
  explanation: string
  evaluation: string
  loading: boolean
}

const DEFAULT_TRACE = {
  trace_id: "abc123",
  spans: [
    {
      span_id: "1",
      trace_id: "abc123",
      name: "auth",
      start_time: new Date().toISOString(),
      end_time: new Date(Date.now() + 120000).toISOString(),
      status: { code: "OK" }
    },
    {
      span_id: "2",
      trace_id: "abc123",
      parent_span_id: "1",
      name: "db",
      start_time: new Date().toISOString(),
      end_time: new Date(Date.now() + 950000).toISOString(),
      status: { code: "ERROR", message: "timeout" }
    },
    {
      span_id: "3",
      trace_id: "abc123",
      parent_span_id: "1",
      name: "api",
      start_time: new Date().toISOString(),
      end_time: new Date(Date.now() + 400000).toISOString(),
      status: { code: "OK" }
    }
  ]
}

function App() {
  const [jsonInput, setJsonInput] = useState(JSON.stringify(DEFAULT_TRACE, null, 2))
  const [researchMode, setResearchMode] = useState(true)

  const [structured, setStructured] = useState<AnalysisState>({ explanation: "", evaluation: "", loading: false })
  const [raw, setRaw] = useState<AnalysisState>({ explanation: "", evaluation: "", loading: false })

  const [facts, setFacts] = useState<SymbolicFact[]>([])
  const [health, setHealth] = useState<SystemHealth | null>(null)
  const [currentTrace, setCurrentTrace] = useState<Trace | null>(null)

  const abortStructured = useRef<AbortController | null>(null)
  const abortRaw = useRef<AbortController | null>(null)

  const streamAnalysis = async (isStructured: boolean, stateSetter: React.Dispatch<React.SetStateAction<AnalysisState>>) => {
    const controller = new AbortController()
    if (isStructured) abortStructured.current = controller
    else abortRaw.current = controller

    stateSetter(prev => ({ ...prev, loading: true, explanation: "", evaluation: "" }))

    try {
      const response = await fetch(`http://localhost:8080/api/analyze?structured=${isStructured}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: jsonInput,
        signal: controller.signal
      })

      if (!response.body) throw new Error('No body')
      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ""

      while (true) {
        const { value, done } = await reader.read()
        if (done) break
        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n\n')
        buffer = lines.pop() || ""

        for (const line of lines) {
          if (line.startsWith('event: metadata')) {
            const data = JSON.parse(line.split('data: ')[1])
            setFacts(data.facts)
            setHealth(data.health)
            setCurrentTrace(data.trace)
          } else if (line.startsWith('event: token')) {
            const token = line.split('data: ')[1]
            stateSetter(prev => ({ ...prev, explanation: prev.explanation + token }))
          }
        }
      }
    } catch (err: any) {
      if (err.name !== 'AbortError') console.error(err)
    } finally {
      stateSetter(prev => ({ ...prev, loading: false }))
    }
  }

  const handleRunAll = () => {
    streamAnalysis(true, setStructured)
    if (researchMode) streamAnalysis(false, setRaw)
  }

  const handleEvaluate = async (isStructured: boolean) => {
    const state = isStructured ? structured : raw
    const setter = isStructured ? setStructured : setRaw

    setter(prev => ({ ...prev, evaluation: "Evaluating performance..." }))

    try {
      const response = await fetch('http://localhost:8080/api/evaluate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          trace: currentTrace,
          facts: facts,
          explanation: state.explanation
        })
      })
      const data = await response.json()
      setter(prev => ({ ...prev, evaluation: data.evaluation }))
    } catch (err) {
      console.error(err)
      setter(prev => ({ ...prev, evaluation: "Evaluation failed." }))
    }
  }

  const getLatency = (span: Span) => {
    return (new Date(span.end_time).getTime() - new Date(span.start_time).getTime())
  }
  const maxLatency = currentTrace?.spans.reduce((max, s) => Math.max(max, getLatency(s)), 0) || 1

  return (
    <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <header>
        <div className="logo">Trace<span>Mind</span><sup>v2.0</sup></div>
        <div className="toggle-container" style={{ marginBottom: 0 }}>
          <label style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <input type="checkbox" checked={researchMode} onChange={() => setResearchMode(!researchMode)} />
            Research Mode <span className="badge purple">BETA</span>
          </label>
        </div>
      </header>

      <main className={`dashboard ${researchMode ? 'research-mode' : ''}`}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <div className="panel" style={{ flex: 1 }}>
            <h3>Telemetry Input</h3>
            <div className="trace-input-area">
              <textarea value={jsonInput} onChange={(e) => setJsonInput(e.target.value)} />
              <button onClick={handleRunAll} disabled={structured.loading || raw.loading}>
                {structured.loading ? 'Processing...' : 'Run Analysis 🚀'}
              </button>
            </div>
          </div>

          {health && (
            <div className="panel" style={{ flex: '0 0 auto', border: '1px solid var(--accent-blue)' }}>
              <h3>System Health Trends</h3>
              <div className="health-stat">
                <span>Error Rate:</span>
                <span style={{ color: health.recent_error_rate > 0.1 ? 'var(--accent-red)' : 'var(--accent-green)' }}>
                  {(health.recent_error_rate * 100).toFixed(1)}%
                </span>
              </div>
              <div className="health-stat">
                <span>Latency Spikes:</span>
                <span style={{ color: 'var(--accent-orange)' }}>{health.slowest_services?.join(', ') || 'None'}</span>
              </div>
            </div>
          )}

          {facts.length > 0 && (
            <div className="panel scrollable" style={{ maxHeight: '300px' }}>
              <h3>Symbolic Facts</h3>
              {facts.map((fact, i) => (
                <div key={i} className={`fact-item ${fact.severity}`}>
                  <div style={{ fontWeight: 'bold', fontSize: '0.7rem' }}>{fact.type}</div>
                  {fact.description}
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="panel research-panel">
          <h3>Structured Reasoning <span className="badge purple">Evolved</span></h3>
          <div className="scrollable">
            <div className="ai-explanation">
              {structured.explanation || (structured.loading ? "Streaming context-aware analysis..." : "Run analysis to see results.")}
            </div>

            {structured.explanation && !structured.loading && (
              <>
                <button className="evaluate-btn" onClick={() => handleEvaluate(true)}>Evaluate with LLM-as-a-Judge</button>
                {structured.evaluation && (
                  <div className="evaluation-result">
                    <strong>SRE Auditor Verdict:</strong><br />
                    {structured.evaluation}
                  </div>
                )}
              </>
            )}

            {currentTrace && (
              <div className="span-visualizer">
                <h5>Trace Breakdown</h5>
                {currentTrace.spans.map((span, i) => (
                  <div key={i} className="span-row">
                    <div className="span-label">{span.name}</div>
                    <div className="span-bar-container">
                      <div className={`span-bar ${span.status.code === 'ERROR' ? 'error' : ''}`} style={{ width: `${(getLatency(span) / maxLatency) * 100}%` }}>
                        {getLatency(span)}ms
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {researchMode && (
          <div className="panel raw-panel">
            <h3>Raw Reasoning <span className="badge">Baseline</span></h3>
            <div className="scrollable">
              <div className="ai-explanation" style={{ color: 'var(--text-secondary)' }}>
                {raw.explanation || (raw.loading ? "Streaming basic analysis..." : "Run analysis to compare.")}
              </div>

              {raw.explanation && !raw.loading && (
                <>
                  <button className="evaluate-btn secondary" onClick={() => handleEvaluate(false)}>Evaluate Baseline</button>
                  {raw.evaluation && (
                    <div className="evaluation-result" style={{ borderColor: 'var(--border-color)', color: 'var(--text-secondary)' }}>
                      <strong>Baseline Audit:</strong><br />
                      {raw.evaluation}
                    </div>
                  )}
                </>
              )}
            </div>
          </div>
        )}
      </main>
    </div>
  )
}

export default App
