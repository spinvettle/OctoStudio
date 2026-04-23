type MessageItem = {
  role: 'user' | 'assistant'
  author: string
  time: string
  content: string[]
}

type ChatMessagesProps = {
  items: MessageItem[]
}

export function ChatMessages({ items }: ChatMessagesProps) {
  return (
    <section className="message-thread" aria-label="Chat messages">
      {items.map((item, index) => (
        <article
          key={`${item.author}-${item.time}-${index}`}
          className={`message-card message-card--${item.role}`}
        >
          <div className="message-card__avatar" aria-hidden="true">
            {item.role === 'assistant' ? 'AI' : 'U'}
          </div>
          <div className="message-card__body">
            <div className="message-card__meta">
              <strong>{item.author}</strong>
              <span>{item.time}</span>
            </div>
            <div className="message-card__content">
              {item.content.map((paragraph) => (
                <p key={paragraph}>{paragraph}</p>
              ))}
            </div>
          </div>
        </article>
      ))}
    </section>
  )
}
