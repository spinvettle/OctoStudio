type ChatComposerProps = {
  hints: string[]
}

export function ChatComposer({ hints }: ChatComposerProps) {
  return (
    <section className="composer-panel">
      <div className="composer-panel__toolbar">
        <button type="button" className="icon-chip">
          Attach
        </button>
        <button type="button" className="icon-chip">
          Search
        </button>
        <button type="button" className="icon-chip">
          Tools
        </button>
      </div>
      <div className="composer-field" aria-label="Prompt editor">
        <p>Ask anything about your project, local models, or deployment workflow...</p>
      </div>
      <div className="composer-panel__footer">
        <div className="hint-list">
          {hints.map((hint) => (
            <span key={hint} className="hint-chip">
              {hint}
            </span>
          ))}
        </div>
        <button type="button" className="send-button">
          Send
        </button>
      </div>
    </section>
  )
}
