export const heroPrimaryAction = {
  label: "Download for Desktop",
  href: "#",
  variant: "primary" as const,
};

export const heroSecondaryAction = {
  label: "Browse Model Library",
  href: "#",
  variant: "secondary" as const,
};

export const heroMetrics = [
  { label: "Concurrent sessions", value: "32" },
  { label: "Average setup", value: "< 5 min" },
  { label: "Inference visibility", value: "Full trace" },
];

export const capabilityGroups = [
  {
    title: "Model lifecycle",
    description:
      "Curate local checkpoints, switch runtimes, and compare variants from the same surface.",
    points: [
      "Versioned model catalog",
      "Quantization-aware metadata",
      "Pinned runtime presets",
    ],
  },
  {
    title: "Prompt and eval loops",
    description:
      "Treat iteration as a product workflow with saved recipes, replayable runs, and clear baselines.",
    points: [
      "Reusable prompt sets",
      "Structured comparison views",
      "Run history by scenario",
    ],
  },
  {
    title: "API handoff",
    description:
      "Move from desktop experimentation to service integration with stable endpoints and familiar schemas.",
    points: [
      "OpenAI-style routes",
      "Streaming responses",
      "Per-project endpoint settings",
    ],
  },
];

export const workflowSteps = [
  {
    step: "01",
    title: "Discover the right checkpoint",
    description:
      "Surface compatibility, size, quantization, and hardware fit before the user commits to a download.",
  },
  {
    step: "02",
    title: "Run locally with full context",
    description:
      "Expose throughput, memory pressure, and session state in the hero visual so the product promise feels operational.",
  },
  {
    step: "03",
    title: "Promote to integration",
    description:
      "Highlight server compatibility and deployment presets to connect exploration with real product use.",
  },
];

export const infrastructureCards = [
  {
    label: "Local inference",
    value: "Metal / CUDA / CPU",
    detail:
      "Hardware-aware execution paths surfaced as a first-class product capability.",
  },
  {
    label: "Server mode",
    value: "OpenAI-compatible",
    detail:
      "A direct bridge from desktop experiments to application integration.",
  },
  {
    label: "Observability",
    value: "Logs, runs, traces",
    detail:
      "Operational language is embedded in the layout instead of hidden in fine print.",
  },
];
