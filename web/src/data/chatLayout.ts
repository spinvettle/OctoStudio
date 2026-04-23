export const sidebarGroups = [
  {
    label: "Workspace",
    items: [
      { label: "Chat", active: true },
      { label: "Models" },
      { label: "Agents", badge: "Beta" },
      { label: "Prompts" },
    ],
  },
  {
    label: "Operations",
    items: [
      { label: "Local Server" },
      { label: "Logs" },
      { label: "Settings" },
    ],
  },
];

export const conversationItems = [
  {
    title: "Landing page rewrite",
    preview: "Rework the product narrative into a chat-first workspace shell.",
    time: "Now",
    active: true,
  },
  {
    title: "Model routing ideas",
    preview: "Compare llama.cpp presets for local coding and reasoning tasks.",
    time: "12m",
  },
  {
    title: "Release checklist",
    preview: "Draft the desktop packaging and API parity checklist.",
    time: "1h",
  },
  {
    title: "Prompt review",
    preview:
      "Tighten system prompts for tool use and multi-step conversations.",
    time: "Yesterday",
  },
];

export const chatHeader = {
  title: "Build a chat-first home screen",
  model: "Qwen3-32B / Local",
  contextWindow: "128K",
  mode: "Balanced",
};

export const messages = [
  {
    role: "assistant" as const,
    author: "OctoStudio",
    time: "09:41",
    content: [
      "Welcome back. Your local workspace is ready, the default coding model is warm, and the API bridge is listening on localhost.",
      "I prepared the latest session summary so you can continue from the product shell instead of a marketing page.",
    ],
  },
  {
    role: "user" as const,
    author: "You",
    time: "09:42",
    content: [
      "Turn the home screen into a modern conversation UI. Keep the product feeling like a serious model workstation.",
    ],
  },
  {
    role: "assistant" as const,
    author: "OctoStudio",
    time: "09:42",
    content: [
      "Done in concept: left navigation for workspace areas, central thread for conversation, model controls in the header, and a persistent composer at the bottom.",
      "This view is display-only for now, so the structure stays clean and ready for real state later.",
    ],
  },
];

export const composerHints = [
  "Summarize this repo",
  "Compare two local models",
  "Draft release notes",
];
