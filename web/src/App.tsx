import './App.css'
import { ChatComposer } from './components/ChatComposer'
import { ChatHeader } from './components/ChatHeader'
import { ChatMessages } from './components/ChatMessages'
import { Sidebar } from './components/Sidebar'
import {
  chatHeader,
  composerHints,
  conversationItems,
  messages,
  sidebarGroups,
} from './data/chatLayout'

function App() {
  return (
    <main className="workspace-shell">
      <Sidebar groups={sidebarGroups} conversations={conversationItems} />
      <section className="chat-workspace">
        <ChatHeader {...chatHeader} />
        <ChatMessages items={messages} />
        <ChatComposer hints={composerHints} />
      </section>
    </main>
  )
}

export default App
