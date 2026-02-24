<template>
  <div class="chat-widget-root">
    <!-- Floating Chat Button -->
    <button
      v-if="!isOpen"
      id="ai-chat-fab"
      class="chat-fab"
      @click="openChat"
      :class="{ 'chat-fab-pulse': showPulse }"
    >
      <div class="chat-fab-inner">
        <i class="ri-robot-2-line"></i>
      </div>
      <span class="chat-fab-label">AI Konsultan</span>
    </button>

    <!-- Chat Panel -->
    <Teleport to="body">
      <Transition name="chat-panel">
        <div v-if="isOpen" class="chat-overlay" @click.self="closeChat">
          <div class="chat-panel">
            <!-- Header -->
            <div class="chat-header">
              <div class="chat-header-left">
                <div class="chat-avatar">
                  <i class="ri-robot-2-fill"></i>
                  <span class="chat-avatar-pulse"></span>
                </div>
                <div class="chat-header-info">
                  <h3>AI Konsultan</h3>
                  <span class="chat-status">
                    <span class="chat-status-dot"></span>
                    Online
                  </span>
                </div>
              </div>
              <div class="chat-header-actions">
                <button class="chat-header-btn" @click="toggleHistory" title="Riwayat">
                  <i class="ri-history-line"></i>
                </button>
                <button class="chat-header-btn" @click="startNewConversation" title="Baru">
                  <i class="ri-add-line"></i>
                </button>
                <button class="chat-header-btn chat-close-btn" @click="closeChat" title="Tutup">
                  <i class="ri-close-line"></i>
                </button>
              </div>
            </div>

            <!-- History Panel -->
            <Transition name="slide-left">
              <div v-if="showHistory" class="chat-history-panel">
                <div class="chat-history-header">
                  <h4>Riwayat Percakapan</h4>
                  <button @click="showHistory = false" class="chat-history-close">
                    <i class="ri-close-line"></i>
                  </button>
                </div>
                <div class="chat-history-list">
                  <div
                    v-for="conv in conversations"
                    :key="conv.id"
                    class="chat-history-item"
                    :class="{ active: currentConversationId === conv.id }"
                    @click="loadConversation(conv.id)"
                  >
                    <div class="chat-history-item-content">
                      <span class="chat-history-title">{{ conv.title }}</span>
                      <span class="chat-history-time">{{ formatTime(conv.updated_at) }}</span>
                    </div>
                    <button
                      class="chat-history-delete"
                      @click.stop="deleteConversation(conv.id)"
                      title="Hapus"
                    >
                      <i class="ri-delete-bin-line"></i>
                    </button>
                  </div>
                  <div v-if="conversations.length === 0" class="chat-history-empty">
                    <i class="ri-chat-new-line"></i>
                    <p>Belum ada percakapan</p>
                  </div>
                </div>
              </div>
            </Transition>

            <!-- Messages -->
            <div class="chat-messages" ref="messagesContainer">
              <!-- Welcome message -->
              <div v-if="messages.length === 0" class="chat-welcome">
                <div class="chat-welcome-icon">
                  <i class="ri-robot-2-fill"></i>
                </div>
                <h3>Halo!</h3>
                <p>Saya AI Konsultan BukaOutlet. Siap membantu seputar kemitraan, bisnis, dan outlet.</p>
                <div class="chat-welcome-actions">
                  <button
                    v-for="action in defaultActions"
                    :key="action.label"
                    class="chat-quick-btn"
                    @click="sendQuickAction(action.action)"
                  >
                    {{ action.label }}
                  </button>
                </div>
              </div>

              <!-- Message list -->
              <div
                v-for="msg in messages"
                :key="msg.id || msg._tempId"
                class="chat-message"
                :class="{ 'chat-message-user': msg.role === 'user', 'chat-message-assistant': msg.role === 'assistant' }"
              >
                <div v-if="msg.role === 'assistant'" class="chat-msg-avatar">
                  <i class="ri-robot-2-fill"></i>
                </div>
                <div class="chat-msg-bubble">
                  <!-- eslint-disable-next-line vue/no-v-html -->
                  <div class="chat-msg-content" v-html="renderMarkdown(msg.content)"></div>
                  <span class="chat-msg-time">{{ formatMsgTime(msg.created_at) }}</span>
                </div>
              </div>

              <!-- Typing indicator -->
              <div v-if="isTyping" class="chat-message chat-message-assistant">
                <div class="chat-msg-avatar">
                  <i class="ri-robot-2-fill"></i>
                </div>
                <div class="chat-msg-bubble chat-typing-bubble">
                  <div class="chat-typing">
                    <span></span><span></span><span></span>
                  </div>
                </div>
              </div>

              <!-- Quick Actions after reply -->
              <div v-if="quickActions.length > 0 && !isTyping" class="chat-quick-actions">
                <button
                  v-for="action in quickActions"
                  :key="action.label"
                  class="chat-quick-btn"
                  @click="sendQuickAction(action.action)"
                >
                  {{ action.label }}
                </button>
              </div>
            </div>

            <!-- Input -->
            <div class="chat-input-area">
              <div class="chat-input-wrapper">
                <textarea
                  ref="chatInput"
                  v-model="inputMessage"
                  @keydown="handleKeydown"
                  placeholder="Ketik pertanyaan..."
                  rows="1"
                  :disabled="isTyping"
                  class="chat-input"
                ></textarea>
                <button
                  class="chat-send-btn"
                  @click="sendMessage"
                  :disabled="!inputMessage.trim() || isTyping"
                >
                  <i class="ri-send-plane-2-fill"></i>
                </button>
              </div>
              <div class="chat-input-hint">AI Konsultan khusus kemitraan &amp; bisnis outlet</div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { chatApi } from '../services/api'

const isOpen = ref(false)
const isTyping = ref(false)
const showHistory = ref(false)
const showPulse = ref(true)
const messages = ref([])
const inputMessage = ref('')
const currentConversationId = ref(null)
const conversations = ref([])
const quickActions = ref([])
const messagesContainer = ref(null)
const chatInput = ref(null)

const defaultActions = [
  { label: 'Cara Jadi Mitra', action: 'Bagaimana cara menjadi mitra?' },
  { label: 'Lihat Outlet', action: 'Outlet apa saja yang tersedia?' },
  { label: 'Cek Budget', action: 'Saya punya budget 50 juta, outlet apa yang cocok?' },
  { label: 'Belajar Bisnis', action: 'Saya ingin belajar bisnis' },
]

let tempIdCounter = 0

function openChat() {
  isOpen.value = true
  showPulse.value = false
  nextTick(() => chatInput.value?.focus())
}

function closeChat() {
  isOpen.value = false
  showHistory.value = false
}

function toggleHistory() {
  showHistory.value = !showHistory.value
  if (showHistory.value) {
    loadConversations()
  }
}

async function loadConversations() {
  try {
    const { data } = await chatApi.conversations()
    conversations.value = data.conversations || []
  } catch {
    conversations.value = []
  }
}

async function loadConversation(convId) {
  try {
    const { data } = await chatApi.messages(convId)
    messages.value = data.messages || []
    currentConversationId.value = convId
    quickActions.value = []
    showHistory.value = false
    scrollToBottom()
  } catch {
    // ignore
  }
}

function startNewConversation() {
  messages.value = []
  currentConversationId.value = null
  quickActions.value = []
  showHistory.value = false
}

async function deleteConversation(convId) {
  try {
    await chatApi.deleteConversation(convId)
    conversations.value = conversations.value.filter(c => c.id !== convId)
    if (currentConversationId.value === convId) {
      startNewConversation()
    }
  } catch {
    // ignore
  }
}

async function sendMessage() {
  const msg = inputMessage.value.trim()
  if (!msg || isTyping.value) return

  const tempId = `temp-${++tempIdCounter}`
  messages.value.push({
    _tempId: tempId,
    role: 'user',
    content: msg,
    created_at: new Date().toISOString(),
  })
  inputMessage.value = ''
  quickActions.value = []
  isTyping.value = true
  scrollToBottom()

  try {
    const { data } = await chatApi.send({
      message: msg,
      conversation_id: currentConversationId.value || '',
    })

    if (data.conversation_id) {
      currentConversationId.value = data.conversation_id
    }

    messages.value.push({
      id: Date.now(),
      role: 'assistant',
      content: data.reply,
      intent: data.intent,
      created_at: new Date().toISOString(),
    })

    quickActions.value = data.quick_actions || []
  } catch {
    messages.value.push({
      id: Date.now(),
      role: 'assistant',
      content: 'Maaf, terjadi kesalahan. Silakan coba lagi.',
      created_at: new Date().toISOString(),
    })
  } finally {
    isTyping.value = false
    scrollToBottom()
    nextTick(() => chatInput.value?.focus())
  }
}

function sendQuickAction(action) {
  inputMessage.value = action
  sendMessage()
}

function handleKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

function renderMarkdown(text) {
  if (!text) return ''
  let html = text
    .replace(/^### (.+)$/gm, '<h4>$1</h4>')
    .replace(/^## (.+)$/gm, '<h3>$1</h3>')
    .replace(/^# (.+)$/gm, '<h2>$1</h2>')
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/_(.+?)_/g, '<em>$1</em>')
    .replace(/```[\s\S]*?```/g, (match) => {
      const code = match.replace(/```\w*\n?/g, '').replace(/```/g, '')
      return `<pre><code>${code}</code></pre>`
    })
    .replace(/`(.+?)`/g, '<code>$1</code>')
    .replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>')
    .replace(/^---$/gm, '<hr>')
    .replace(/^- (.+)$/gm, '<li>$1</li>')
    .replace(/^\d+\.\s(.+)$/gm, '<li>$1</li>')
    .replace(/\n\n/g, '<br>')
    .replace(/\n/g, '<br>')

  html = html.replace(/((?:<li>.*?<\/li>\s*(?:<br>)?)+)/g, '<ul>$1</ul>')
  html = html.replace(/<ul>([\s\S]*?)<\/ul>/g, (match) => {
    return match.replace(/<br>/g, '')
  })

  return html
}

function formatTime(date) {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now - d

  if (diff < 60000) return 'Baru saja'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m lalu`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}j lalu`
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

function formatMsgTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  setTimeout(() => {
    showPulse.value = true
  }, 3000)
})
</script>

<style scoped>
/* ══════════════════════════════════════
   ROOT
   ══════════════════════════════════════ */
.chat-widget-root {
  position: relative;
  z-index: 9998;
}

/* ══════════════════════════════════════
   FLOATING ACTION BUTTON
   ══════════════════════════════════════ */
.chat-fab {
  position: fixed !important;
  bottom: 20px;
  right: 16px;
  z-index: 9998;
  display: flex !important;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
  border: none !important;
  border-radius: 28px;
  padding: 12px 18px 12px 12px;
  cursor: pointer;
  box-shadow: 0 6px 24px rgba(102, 126, 234, 0.45), 0 2px 8px rgba(0,0,0,0.15);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-family: inherit;
}

.chat-fab:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 32px rgba(102, 126, 234, 0.55);
}

.chat-fab-inner {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255,255,255,0.2) !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: white !important;
}

.chat-fab-label {
  font-size: 13px;
  font-weight: 600;
  color: white !important;
  letter-spacing: -0.01em;
}

.chat-fab-pulse {
  animation: fabPulse 2s ease-in-out infinite;
}

@keyframes fabPulse {
  0%, 100% { box-shadow: 0 6px 24px rgba(102,126,234,0.45); }
  50% { box-shadow: 0 6px 24px rgba(102,126,234,0.45), 0 0 0 10px rgba(102,126,234,0.12); }
}

/* ══════════════════════════════════════
   CHAT OVERLAY & PANEL (Mobile Full Screen)
   ══════════════════════════════════════ */
.chat-overlay {
  position: fixed !important;
  inset: 0;
  z-index: 9999;
  display: flex !important;
  background: rgba(0,0,0,0.4) !important;
  backdrop-filter: blur(2px);
}

.chat-panel {
  width: 100% !important;
  height: 100% !important;
  background: #fff !important;
  display: flex !important;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

/* ══════════════════════════════════════
   HEADER
   ══════════════════════════════════════ */
.chat-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
  padding: 14px 16px;
  padding-top: calc(14px + env(safe-area-inset-top, 0px));
  display: flex !important;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.chat-header-left {
  display: flex !important;
  align-items: center;
  gap: 10px;
}

.chat-avatar {
  width: 38px;
  height: 38px;
  border-radius: 12px;
  background: rgba(255,255,255,0.2) !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: white !important;
  position: relative;
}

.chat-avatar-pulse {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #4ade80 !important;
  border: 2px solid #667eea !important;
}

.chat-header-info h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 700;
  color: white !important;
  letter-spacing: -0.01em;
}

.chat-status {
  display: flex !important;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: rgba(255,255,255,0.75) !important;
}

.chat-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #4ade80 !important;
  animation: statusPulse 2s infinite;
}

@keyframes statusPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.chat-header-actions {
  display: flex !important;
  gap: 6px;
}

.chat-header-btn {
  width: 36px;
  height: 36px;
  border: none !important;
  border-radius: 10px;
  background: rgba(255,255,255,0.15) !important;
  color: white !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 18px;
  transition: background 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.chat-header-btn:hover,
.chat-header-btn:active {
  background: rgba(255,255,255,0.3) !important;
}

.chat-close-btn:hover,
.chat-close-btn:active {
  background: rgba(239,68,68,0.5) !important;
}

/* ══════════════════════════════════════
   HISTORY PANEL
   ══════════════════════════════════════ */
.chat-history-panel {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #f8f9fb !important;
  z-index: 10;
  display: flex !important;
  flex-direction: column;
}

.chat-history-header {
  display: flex !important;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  padding-top: calc(16px + env(safe-area-inset-top, 0px));
  border-bottom: 1px solid #e5e7eb !important;
  background: white !important;
}

.chat-history-header h4 {
  margin: 0;
  font-size: 15px;
  font-weight: 700;
  color: #1f2937 !important;
}

.chat-history-close {
  width: 34px;
  height: 34px;
  border: none !important;
  border-radius: 10px;
  background: #f1f5f9 !important;
  color: #64748b !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 18px;
}

.chat-history-list {
  overflow-y: auto;
  flex: 1;
  padding: 8px;
  -webkit-overflow-scrolling: touch;
}

.chat-history-item {
  display: flex !important;
  align-items: center;
  justify-content: space-between;
  padding: 14px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
  margin-bottom: 4px;
  -webkit-tap-highlight-color: transparent;
}

.chat-history-item:active {
  background: #eef0f4 !important;
}

.chat-history-item.active {
  background: #e8ebf5 !important;
  border-left: 3px solid #667eea !important;
}

.chat-history-item-content {
  flex: 1;
  min-width: 0;
}

.chat-history-title {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151 !important;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chat-history-time {
  display: block;
  font-size: 11px;
  color: #9ca3af !important;
  margin-top: 3px;
}

.chat-history-delete {
  width: 34px;
  height: 34px;
  border: none !important;
  border-radius: 8px;
  background: transparent !important;
  color: #9ca3af !important;
  cursor: pointer;
  display: flex !important;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
  -webkit-tap-highlight-color: transparent;
}

.chat-history-delete:active {
  background: #fee2e2 !important;
  color: #ef4444 !important;
}

.chat-history-empty {
  display: flex !important;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 20px;
  color: #9ca3af !important;
  text-align: center;
}

.chat-history-empty i {
  font-size: 36px;
  margin-bottom: 12px;
  opacity: 0.5;
  color: #9ca3af !important;
}

.chat-history-empty p {
  margin: 0;
  font-size: 13px;
  color: #9ca3af !important;
}

/* ══════════════════════════════════════
   MESSAGES AREA
   ══════════════════════════════════════ */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex !important;
  flex-direction: column;
  gap: 12px;
  background: #f8f9fb !important;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
}

.chat-messages::-webkit-scrollbar { width: 0; }

/* Welcome */
.chat-welcome {
  display: flex !important;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 24px 16px;
  gap: 8px;
}

.chat-welcome-icon {
  width: 56px;
  height: 56px;
  border-radius: 18px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: white !important;
  margin-bottom: 4px;
  box-shadow: 0 6px 20px rgba(102,126,234,0.3);
}

.chat-welcome h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #1f2937 !important;
}

.chat-welcome p {
  margin: 0;
  font-size: 13px;
  color: #6b7280 !important;
  line-height: 1.5;
  max-width: 280px;
}

.chat-welcome-actions {
  display: flex !important;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: center;
  margin-top: 8px;
}

/* Messages */
.chat-message {
  display: flex !important;
  gap: 8px;
  animation: messageIn 0.3s ease-out;
}

@keyframes messageIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.chat-message-user {
  flex-direction: row-reverse;
}

.chat-msg-avatar {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  color: white !important;
  font-size: 14px;
  flex-shrink: 0;
  margin-top: 2px;
}

.chat-msg-bubble {
  max-width: 85%;
  padding: 10px 14px;
  border-radius: 16px;
  position: relative;
}

.chat-message-user .chat-msg-bubble {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
  border-bottom-right-radius: 4px;
}

.chat-message-assistant .chat-msg-bubble {
  background: white !important;
  color: #1f2937 !important;
  border-bottom-left-radius: 4px;
  box-shadow: 0 1px 6px rgba(0,0,0,0.06);
}

.chat-msg-content {
  font-size: 14px;
  line-height: 1.6;
  word-break: break-word;
  color: inherit !important;
}

.chat-msg-content :deep(h2) {
  font-size: 15px;
  font-weight: 700;
  margin: 6px 0;
  color: #1f2937 !important;
}

.chat-msg-content :deep(h3) {
  font-size: 14px;
  font-weight: 700;
  margin: 6px 0 4px;
  color: #374151 !important;
}

.chat-msg-content :deep(h4) {
  font-size: 13px;
  font-weight: 600;
  margin: 4px 0 3px;
  color: #4b5563 !important;
}

.chat-msg-content :deep(ul) {
  margin: 4px 0;
  padding-left: 10px;
  list-style: none;
}

.chat-msg-content :deep(li) {
  position: relative;
  margin-bottom: 3px;
  font-size: 13px;
  line-height: 1.5;
  color: inherit !important;
}

.chat-msg-content :deep(strong) {
  font-weight: 600;
  color: inherit !important;
}

.chat-msg-content :deep(code) {
  background: rgba(102,126,234,0.08) !important;
  padding: 2px 5px;
  border-radius: 4px;
  font-family: 'SF Mono', monospace;
  font-size: 12px;
  color: inherit !important;
}

.chat-msg-content :deep(pre) {
  background: #1e293b !important;
  color: #e2e8f0 !important;
  padding: 10px;
  border-radius: 8px;
  overflow-x: auto;
  font-size: 12px;
  margin: 6px 0;
}

.chat-msg-content :deep(blockquote) {
  border-left: 3px solid #667eea !important;
  padding: 6px 10px;
  margin: 6px 0;
  background: rgba(102,126,234,0.06) !important;
  border-radius: 0 8px 8px 0;
  font-size: 13px;
  color: #4b5563 !important;
}

.chat-msg-content :deep(hr) {
  border: none !important;
  border-top: 1px solid #e5e7eb !important;
  margin: 8px 0;
}

.chat-msg-time {
  display: block;
  font-size: 10px;
  opacity: 0.5;
  margin-top: 4px;
  text-align: right;
  color: inherit !important;
}

/* Typing indicator */
.chat-typing-bubble {
  padding: 12px 18px;
}

.chat-typing {
  display: flex !important;
  gap: 4px;
  align-items: center;
}

.chat-typing span {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #9ca3af !important;
  animation: typingBounce 1.4s infinite;
}

.chat-typing span:nth-child(2) { animation-delay: 0.2s; }
.chat-typing span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typingBounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
  30% { transform: translateY(-5px); opacity: 1; }
}

/* Quick Actions */
.chat-quick-actions {
  display: flex !important;
  flex-wrap: wrap;
  gap: 6px;
  padding: 4px 0;
}

.chat-quick-btn {
  background: white !important;
  border: 1.5px solid #e5e7eb !important;
  border-radius: 20px;
  padding: 8px 14px;
  font-size: 13px;
  color: #374151 !important;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  -webkit-tap-highlight-color: transparent;
}

.chat-quick-btn:active {
  background: #f0f1ff !important;
  border-color: #667eea !important;
  color: #667eea !important;
}

/* ══════════════════════════════════════
   INPUT AREA
   ══════════════════════════════════════ */
.chat-input-area {
  padding: 12px 16px;
  padding-bottom: calc(12px + env(safe-area-inset-bottom, 0px));
  background: white !important;
  border-top: 1px solid #f0f1f3 !important;
  flex-shrink: 0;
}

.chat-input-wrapper {
  display: flex !important;
  align-items: flex-end;
  gap: 8px;
  background: #f3f4f6 !important;
  border-radius: 14px;
  padding: 6px 6px 6px 14px;
  transition: all 0.2s;
  border: 1.5px solid transparent !important;
}

.chat-input-wrapper:focus-within {
  background: white !important;
  border-color: #667eea !important;
  box-shadow: 0 0 0 3px rgba(102,126,234,0.1);
}

.chat-input {
  flex: 1;
  border: none !important;
  background: transparent !important;
  font-size: 15px;
  color: #1f2937 !important;
  resize: none;
  outline: none !important;
  font-family: inherit;
  max-height: 100px;
  line-height: 1.5;
  padding: 6px 0;
  -webkit-appearance: none;
}

.chat-input::placeholder {
  color: #9ca3af !important;
}

.chat-send-btn {
  width: 36px;
  height: 36px;
  border: none !important;
  border-radius: 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  font-size: 17px;
  cursor: pointer;
  flex-shrink: 0;
  transition: all 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.chat-send-btn:active:not(:disabled) {
  transform: scale(0.95);
}

.chat-send-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.chat-input-hint {
  text-align: center;
  margin-top: 6px;
  font-size: 10px;
  color: #b0b7c3 !important;
}

/* ══════════════════════════════════════
   TRANSITIONS
   ══════════════════════════════════════ */
.chat-panel-enter-active { transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.chat-panel-leave-active { transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); }
.chat-panel-enter-from { opacity: 0; }
.chat-panel-leave-to { opacity: 0; }
.chat-panel-enter-from .chat-panel { transform: translateY(100%); }
.chat-panel-leave-to .chat-panel { transform: translateY(100%); }

.slide-left-enter-active { transition: all 0.3s ease; }
.slide-left-leave-active { transition: all 0.2s ease; }
.slide-left-enter-from { transform: translateX(-100%); opacity: 0; }
.slide-left-leave-to { transform: translateX(-100%); opacity: 0; }
</style>
