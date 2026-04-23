type FeatureItem = {
  title: string
  description: string
  points: string[]
}

type FeatureGridProps = {
  items: FeatureItem[]
}

export function FeatureGrid({ items }: FeatureGridProps) {
  return (
    <div className="feature-grid">
      {items.map((item) => (
        <article key={item.title} className="surface-card">
          <h3>{item.title}</h3>
          <p>{item.description}</p>
          <ul className="surface-card__list">
            {item.points.map((point) => (
              <li key={point}>{point}</li>
            ))}
          </ul>
        </article>
      ))}
    </div>
  )
}
