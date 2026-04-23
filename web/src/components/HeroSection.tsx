type Action = {
  label: string
  href: string
  variant: 'primary' | 'secondary'
}

type Metric = {
  value: string
  label: string
}

type HeroSectionProps = {
  eyebrow: string
  title: string
  description: string
  primaryAction: Action
  secondaryAction: Action
  metrics: Metric[]
}

function HeroVisual() {
  return (
    <div className="hero-visual" aria-hidden="true">
      <div className="hero-visual__frame hero-visual__frame--main">
        <div className="hero-visual__topbar">
          <span />
          <span />
          <span />
        </div>
        <div className="hero-visual__content">
          <div className="hero-visual__sidebar">
            <div />
            <div />
            <div />
          </div>
          <div className="hero-visual__panel">
            <div className="hero-visual__model-row">
              <strong>llama-3.3-70b</strong>
              <span>Loaded on Metal</span>
            </div>
            <div className="hero-visual__chart">
              <span />
              <span />
              <span />
              <span />
              <span />
              <span />
            </div>
            <div className="hero-visual__terminal">
              <span>localhost:1234/v1/chat/completions</span>
              <span>18.4 tok/s</span>
              <span>Ready</span>
            </div>
          </div>
        </div>
      </div>
      <div className="hero-visual__frame hero-visual__frame--floating">
        <p>Deploy preset</p>
        <strong>OpenAI-compatible API</strong>
        <span>Streaming enabled</span>
      </div>
    </div>
  )
}

export function HeroSection({
  eyebrow,
  title,
  description,
  primaryAction,
  secondaryAction,
  metrics,
}: HeroSectionProps) {
  const actions = [primaryAction, secondaryAction]

  return (
    <section className="hero-section">
      <div className="hero-copy">
        <p className="hero-copy__eyebrow">{eyebrow}</p>
        <h1>{title}</h1>
        <p className="hero-copy__description">{description}</p>
        <div className="hero-copy__actions">
          {actions.map((action) => (
            <a
              key={action.label}
              className={`button button--${action.variant}`}
              href={action.href}
            >
              {action.label}
            </a>
          ))}
        </div>
        <dl className="hero-copy__metrics">
          {metrics.map((metric) => (
            <div key={metric.label} className="metric-card">
              <dt>{metric.label}</dt>
              <dd>{metric.value}</dd>
            </div>
          ))}
        </dl>
      </div>
      <HeroVisual />
    </section>
  )
}
