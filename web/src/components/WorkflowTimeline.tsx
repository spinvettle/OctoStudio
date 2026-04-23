type WorkflowStep = {
  step: string
  title: string
  description: string
}

type WorkflowTimelineProps = {
  items: WorkflowStep[]
}

export function WorkflowTimeline({ items }: WorkflowTimelineProps) {
  return (
    <ol className="workflow-timeline">
      {items.map((item) => (
        <li key={item.step} className="workflow-timeline__item">
          <p className="workflow-timeline__step">{item.step}</p>
          <div>
            <h3>{item.title}</h3>
            <p>{item.description}</p>
          </div>
        </li>
      ))}
    </ol>
  )
}
