type ChatHeaderProps = {
  title: string
  model: string
  contextWindow: string
  mode: string
}

export function ChatHeader({
  title,
  model,
  contextWindow,
  mode,
}: ChatHeaderProps) {
  return (
    <header className="chat-header">
      <div>
        <p className="chat-header__eyebrow">Conversation</p>
        <h1>{title}</h1>
      </div>
      <div className="chat-header__controls">
        <div className="pill-control">
          <span className="pill-control__label">Model</span>
          <strong>{model}</strong>
        </div>
        <div className="pill-control">
          <span className="pill-control__label">Context</span>
          <strong>{contextWindow}</strong>
        </div>
        <div className="pill-control">
          <span className="pill-control__label">Mode</span>
          <strong>{mode}</strong>
        </div>
      </div>
    </header>
  )
}
