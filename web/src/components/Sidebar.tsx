type SidebarGroup = {
  label: string
  items: Array<{
    label: string
    active?: boolean
    badge?: string
  }>
}

type ConversationItem = {
  title: string
  preview: string
  time: string
  active?: boolean
}

type SidebarProps = {
  groups: SidebarGroup[]
  conversations: ConversationItem[]
}

export function Sidebar({ groups, conversations }: SidebarProps) {
  return (
    <aside className="sidebar">
      <div className="brand-panel">
        <div className="brand-mark" aria-hidden="true">
          OS
        </div>
        <div>
          <strong>OctoStudio</strong>
          <p>Local model workspace</p>
        </div>
      </div>

      {groups.map((group) => (
        <section key={group.label} className="sidebar-section">
          <p className="sidebar-section__label">{group.label}</p>
          <div className="sidebar-nav">
            {group.items.map((item) => (
              <button
                key={item.label}
                type="button"
                className={`sidebar-nav__item${item.active ? ' is-active' : ''}`}
              >
                <span>{item.label}</span>
                {item.badge ? <span className="sidebar-nav__badge">{item.badge}</span> : null}
              </button>
            ))}
          </div>
        </section>
      ))}

      <section className="sidebar-section sidebar-section--grow">
        <div className="sidebar-section__row">
          <p className="sidebar-section__label">Recent chats</p>
          <button type="button" className="icon-button" aria-label="New chat">
            +
          </button>
        </div>
        <div className="conversation-list">
          {conversations.map((conversation) => (
            <article
              key={`${conversation.title}-${conversation.time}`}
              className={`conversation-card${conversation.active ? ' is-active' : ''}`}
            >
              <div className="conversation-card__meta">
                <strong>{conversation.title}</strong>
                <span>{conversation.time}</span>
              </div>
              <p>{conversation.preview}</p>
            </article>
          ))}
        </div>
      </section>
    </aside>
  )
}
