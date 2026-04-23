type MetricStripItem = {
  label: string
  value: string
  detail: string
}

type MetricStripProps = {
  items: MetricStripItem[]
}

export function MetricStrip({ items }: MetricStripProps) {
  return (
    <div className="metric-strip">
      {items.map((item) => (
        <article key={item.label} className="surface-card surface-card--compact">
          <p className="metric-strip__label">{item.label}</p>
          <strong>{item.value}</strong>
          <p>{item.detail}</p>
        </article>
      ))}
    </div>
  )
}
