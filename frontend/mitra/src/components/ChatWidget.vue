<template>
  <div class="cw-root">
    <!-- Floating Button -->
    <button v-if="!isOpen" id="ai-chat-fab" class="cw-fab" @click="openChat" :class="{ pulse: showPulse }">
      <span class="cw-fab-icon"><i class="ri-robot-2-line"></i></span>
      <span class="cw-fab-text">AI Konsultan</span>
    </button>

    <!-- Chat Panel -->
    <Teleport to="body">
      <Transition name="cw-slide">
        <div v-if="isOpen" class="cw-overlay" @click.self="closeChat">
          <div class="cw-panel">
            <!-- Header -->
            <header class="cw-header">
              <div class="cw-header-left">
                <div class="cw-logo">
                  <i class="ri-robot-2-fill"></i>
                  <span class="cw-online-dot"></span>
                </div>
                <div>
                  <h3 class="cw-title">AI Konsultan</h3>
                  <span class="cw-subtitle"><span class="cw-status-dot"></span> Online — siap membantu</span>
                </div>
              </div>
              <div class="cw-header-right">
                <button class="cw-hbtn" @click="toggleHistory" title="Riwayat"><i class="ri-history-line"></i></button>
                <button class="cw-hbtn" @click="startNewConversation" title="Baru"><i class="ri-add-line"></i></button>
                <button class="cw-hbtn cw-hbtn-close" @click="closeChat" title="Tutup"><i class="ri-close-line"></i></button>
              </div>
            </header>

            <!-- History Drawer -->
            <Transition name="cw-drawer">
              <div v-if="showHistory" class="cw-history">
                <div class="cw-history-top">
                  <h4>Riwayat Percakapan</h4>
                  <button class="cw-hbtn" @click="showHistory = false"><i class="ri-close-line"></i></button>
                </div>
                <div class="cw-history-list">
                  <div v-for="conv in conversations" :key="conv.id"
                    class="cw-history-item" :class="{ active: currentConversationId === conv.id }"
                    @click="loadConversation(conv.id)">
                    <div class="cw-history-body">
                      <span class="cw-history-title">{{ conv.title }}</span>
                      <span class="cw-history-time">{{ formatTime(conv.updated_at) }}</span>
                    </div>
                    <button class="cw-history-del" @click.stop="deleteConversation(conv.id)"><i class="ri-delete-bin-6-line"></i></button>
                  </div>
                  <div v-if="!conversations.length" class="cw-history-empty">
                    <i class="ri-chat-new-line"></i>
                    <p>Belum ada percakapan</p>
                  </div>
                </div>
              </div>
            </Transition>

            <!-- Messages -->
            <div class="cw-messages" ref="messagesContainer">
              <!-- Welcome -->
              <div v-if="!messages.length" class="cw-welcome">
                <div class="cw-welcome-icon"><i class="ri-robot-2-fill"></i></div>
                <h3>Halo! 👋</h3>
                <p>Saya AI Konsultan BukaOutlet — siap bantu seputar franchise, kemitraan, dan bisnis outlet.</p>
                <div class="cw-welcome-grid">
                  <button v-for="a in defaultActions" :key="a.label" class="cw-chip" @click="sendQuickAction(a.action)">
                    <i :class="a.icon"></i> {{ a.label }}
                  </button>
                </div>
              </div>

              <!-- Messages -->
              <div v-for="msg in messages" :key="msg.id || msg._tempId"
                class="cw-msg" :class="msg.role === 'user' ? 'cw-msg-user' : 'cw-msg-ai'">
                <div v-if="msg.role === 'assistant'" class="cw-msg-ava"><i class="ri-robot-2-fill"></i></div>
                <div class="cw-bubble">
                  <!-- eslint-disable-next-line vue/no-v-html -->
                  <div class="cw-prose" v-html="renderMarkdown(msg.content)"></div>
                  <time class="cw-time">{{ formatMsgTime(msg.created_at) }}</time>
                </div>
              </div>

              <!-- Typing -->
              <div v-if="isTyping" class="cw-msg cw-msg-ai">
                <div class="cw-msg-ava"><i class="ri-robot-2-fill"></i></div>
                <div class="cw-bubble cw-bubble-typing">
                  <div class="cw-dots"><span></span><span></span><span></span></div>
                  <span class="cw-typing-text">Sedang mencari info...</span>
                </div>
              </div>

              <!-- Quick Actions -->
              <div v-if="quickActions.length && !isTyping" class="cw-actions">
                <button v-for="a in quickActions" :key="a.label" class="cw-chip" @click="sendQuickAction(a.action)">
                  {{ a.label }}
                </button>
              </div>
            </div>

            <!-- Input -->
            <div class="cw-input-area">
              <div class="cw-input-box" :class="{ focused: inputFocused }">
                <textarea ref="chatInput" v-model="inputMessage" @keydown="handleKeydown"
                  @focus="inputFocused = true" @blur="inputFocused = false"
                  placeholder="Tanya tentang franchise, bisnis, outlet..." rows="1"
                  :disabled="isTyping" class="cw-textarea"></textarea>
                <button class="cw-send" @click="sendMessage" :disabled="!inputMessage.trim() || isTyping">
                  <i class="ri-send-plane-2-fill"></i>
                </button>
              </div>
              <p class="cw-disclaimer">Didukung AI — jawaban bisa mencakup info franchise umum & BukaOutlet</p>
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
const inputFocused = ref(false)
const messages = ref([])
const inputMessage = ref('')
const currentConversationId = ref(null)
const conversations = ref([])
const quickActions = ref([])
const messagesContainer = ref(null)
const chatInput = ref(null)

const defaultActions = [
  { label: 'Franchise Viral', action: 'Franchise apa yang lagi viral di Indonesia?', icon: 'ri-fire-line' },
  { label: 'Cara Jadi Mitra', action: 'Bagaimana cara menjadi mitra BukaOutlet?', icon: 'ri-user-add-line' },
  { label: 'Cek Budget', action: 'Rekomendasi franchise modal di bawah 50 juta', icon: 'ri-money-dollar-circle-line' },
  { label: 'Tips Bisnis', action: 'Tips memulai bisnis franchise untuk pemula', icon: 'ri-lightbulb-line' },
]

let tempIdCounter = 0

function openChat() {
  isOpen.value = true
  showPulse.value = false
  nextTick(() => chatInput.value?.focus())
}
function closeChat() { isOpen.value = false; showHistory.value = false }
function toggleHistory() {
  showHistory.value = !showHistory.value
  if (showHistory.value) loadConversations()
}
async function loadConversations() {
  try { const { data } = await chatApi.conversations(); conversations.value = data.conversations || [] }
  catch { conversations.value = [] }
}
async function loadConversation(convId) {
  try {
    const { data } = await chatApi.messages(convId)
    messages.value = data.messages || []
    currentConversationId.value = convId
    quickActions.value = []
    showHistory.value = false
    scrollToBottom()
  } catch { /* ignore */ }
}
function startNewConversation() {
  messages.value = []; currentConversationId.value = null; quickActions.value = []; showHistory.value = false
}
async function deleteConversation(convId) {
  try {
    await chatApi.deleteConversation(convId)
    conversations.value = conversations.value.filter(c => c.id !== convId)
    if (currentConversationId.value === convId) startNewConversation()
  } catch { /* ignore */ }
}
async function sendMessage() {
  const msg = inputMessage.value.trim()
  if (!msg || isTyping.value) return
  messages.value.push({ _tempId: `t${++tempIdCounter}`, role: 'user', content: msg, created_at: new Date().toISOString() })
  inputMessage.value = ''
  quickActions.value = []
  isTyping.value = true
  scrollToBottom()
  try {
    const { data } = await chatApi.send({ message: msg, conversation_id: currentConversationId.value || '' })
    if (data.conversation_id) currentConversationId.value = data.conversation_id
    messages.value.push({ id: Date.now(), role: 'assistant', content: data.reply, created_at: new Date().toISOString() })
    quickActions.value = data.quick_actions || []
  } catch {
    messages.value.push({ id: Date.now(), role: 'assistant', content: 'Maaf, terjadi kesalahan. Silakan coba lagi.', created_at: new Date().toISOString() })
  } finally { isTyping.value = false; scrollToBottom(); nextTick(() => chatInput.value?.focus()) }
}
function sendQuickAction(a) { inputMessage.value = a; sendMessage() }
function handleKeydown(e) { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage() } }
function scrollToBottom() { nextTick(() => { if (messagesContainer.value) messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight }) }

/* ── Proper Markdown Renderer for GPT output ── */
function renderMarkdown(text) {
  if (!text) return ''
  let s = text

  // Escape HTML
  s = s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  // Code blocks (``` ... ```)
  s = s.replace(/```(\w*)\n([\s\S]*?)```/g, (_, lang, code) => {
    return `<pre class="cw-code"><code>${code.trim()}</code></pre>`
  })

  // Inline code
  s = s.replace(/`([^`]+)`/g, '<code class="cw-inline-code">$1</code>')

  // Headings
  s = s.replace(/^#### (.+)$/gm, '<h5 class="cw-h5">$1</h5>')
  s = s.replace(/^### (.+)$/gm, '<h4 class="cw-h4">$1</h4>')
  s = s.replace(/^## (.+)$/gm, '<h3 class="cw-h3">$1</h3>')
  s = s.replace(/^# (.+)$/gm, '<h2 class="cw-h2">$1</h2>')

  // Bold & italic
  s = s.replace(/\*\*\*(.+?)\*\*\*/g, '<strong><em>$1</em></strong>')
  s = s.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
  s = s.replace(/(?<!\*)\*([^*]+?)\*(?!\*)/g, '<em>$1</em>')

  // Horizontal rule
  s = s.replace(/^---$/gm, '<hr class="cw-hr">')

  // Blockquote
  s = s.replace(/^&gt; (.+)$/gm, '<blockquote class="cw-bq">$1</blockquote>')

  // Links  [text](url)
  s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener" class="cw-link">$1</a>')

  // Tables
  s = s.replace(/^(\|.+\|)\n(\|[\s-:|]+\|)\n((?:\|.+\|\n?)+)/gm, (_, header, sep, body) => {
    const thCells = header.split('|').filter(c => c.trim()).map(c => `<th>${c.trim()}</th>`).join('')
    const rows = body.trim().split('\n').map(row => {
      const cells = row.split('|').filter(c => c.trim()).map(c => `<td>${c.trim()}</td>`).join('')
      return `<tr>${cells}</tr>`
    }).join('')
    return `<table class="cw-table"><thead><tr>${thCells}</tr></thead><tbody>${rows}</tbody></table>`
  })

  // Process lines for lists and paragraphs
  const lines = s.split('\n')
  let html = ''
  let inUl = false, inOl = false

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const ulMatch = line.match(/^[\s]*[-•] (.+)/)
    const olMatch = line.match(/^[\s]*(\d+)\.\s(.+)/)

    if (ulMatch) {
      if (!inUl) { html += '<ul class="cw-ul">'; inUl = true }
      if (inOl) { html += '</ol>'; inOl = false }
      html += `<li>${ulMatch[1]}</li>`
    } else if (olMatch) {
      if (!inOl) { html += '<ol class="cw-ol">'; inOl = true }
      if (inUl) { html += '</ul>'; inUl = false }
      html += `<li>${olMatch[2]}</li>`
    } else {
      if (inUl) { html += '</ul>'; inUl = false }
      if (inOl) { html += '</ol>'; inOl = false }
      if (line.trim() === '') {
        html += '<div class="cw-spacer"></div>'
      } else if (line.startsWith('<')) {
        html += line
      } else {
        html += `<p class="cw-p">${line}</p>`
      }
    }
  }
  if (inUl) html += '</ul>'
  if (inOl) html += '</ol>'

  return html
}

function formatTime(date) {
  if (!date) return ''
  const d = new Date(date), now = new Date(), diff = now - d
  if (diff < 60000) return 'Baru saja'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m lalu`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}j lalu`
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}
function formatMsgTime(date) {
  if (!date) return ''
  return new Date(date).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => { setTimeout(() => { showPulse.value = true }, 3000) })
</script>

<style scoped>
/* ════════════════════════════════════════════
   RUSSIAN ORANGE PALETTE
   Primary: #E8702A  Darker: #C45A1C  Lighter: #FFF3EC
   Gradient: #E8702A → #D4451A
   ════════════════════════════════════════════ */
:root {
  --cw-primary: #E8702A;
  --cw-primary-dark: #C45A1C;
  --cw-primary-light: #FFF3EC;
  --cw-primary-glow: rgba(232, 112, 42, 0.35);
  --cw-gradient: linear-gradient(135deg, #E8702A 0%, #D4451A 100%);
  --cw-text: #1A1A2E;
  --cw-text-secondary: #6B7280;
  --cw-bg: #FAFAFA;
  --cw-card: #FFFFFF;
  --cw-border: #F0F0F0;
  --cw-radius: 16px;
}

/* ═══ ROOT ═══ */
.cw-root { position: relative; z-index: 9998; }

/* ═══ FAB BUTTON ═══ */
.cw-fab {
  position: fixed !important; bottom: 20px; right: 16px; z-index: 9998;
  display: flex !important; align-items: center; gap: 10px;
  background: var(--cw-gradient) !important; color: #fff !important;
  border: none !important; border-radius: 28px;
  padding: 12px 20px 12px 14px; cursor: pointer;
  box-shadow: 0 8px 28px var(--cw-primary-glow), 0 2px 8px rgba(0,0,0,0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-family: 'Inter', system-ui, sans-serif;
}
.cw-fab:hover { transform: translateY(-3px); box-shadow: 0 12px 36px var(--cw-primary-glow); }
.cw-fab-icon {
  width: 34px; height: 34px; border-radius: 50%;
  background: rgba(255,255,255,0.2) !important;
  display: flex !important; align-items: center; justify-content: center;
  font-size: 19px; color: #fff !important;
}
.cw-fab-text { font-size: 13px; font-weight: 600; color: #fff !important; }
.cw-fab.pulse { animation: cwPulse 2.5s ease-in-out infinite; }
@keyframes cwPulse {
  0%,100% { box-shadow: 0 8px 28px var(--cw-primary-glow); }
  50% { box-shadow: 0 8px 28px var(--cw-primary-glow), 0 0 0 12px rgba(232,112,42,0.1); }
}

/* ═══ OVERLAY & PANEL ═══ */
.cw-overlay {
  position: fixed !important; inset: 0; z-index: 9999;
  display: flex !important; background: rgba(0,0,0,0.45) !important;
  backdrop-filter: blur(4px); -webkit-backdrop-filter: blur(4px);
}
.cw-panel {
  width: 100% !important; height: 100% !important;
  background: var(--cw-bg) !important;
  display: flex !important; flex-direction: column;
  overflow: hidden; position: relative;
}

/* Desktop: right-side panel */
@media (min-width: 640px) {
  .cw-panel {
    width: 420px !important; height: 700px !important; max-height: 90vh;
    border-radius: 20px; margin: auto;
    box-shadow: 0 25px 60px rgba(0,0,0,0.2);
  }
}

/* ═══ HEADER ═══ */
.cw-header {
  background: var(--cw-gradient) !important; color: #fff !important;
  padding: 16px 16px; padding-top: calc(16px + env(safe-area-inset-top, 0px));
  display: flex !important; align-items: center; justify-content: space-between;
  flex-shrink: 0;
}
.cw-header-left { display: flex !important; align-items: center; gap: 12px; }
.cw-logo {
  width: 42px; height: 42px; border-radius: 14px;
  background: rgba(255,255,255,0.18) !important;
  display: flex !important; align-items: center; justify-content: center;
  font-size: 22px; color: #fff !important; position: relative;
}
.cw-online-dot {
  position: absolute; bottom: -1px; right: -1px;
  width: 11px; height: 11px; border-radius: 50%;
  background: #34D399 !important; border: 2.5px solid var(--cw-primary) !important;
}
.cw-title { margin: 0; font-size: 16px; font-weight: 700; color: #fff !important; letter-spacing: -0.02em; }
.cw-subtitle { font-size: 11px; color: rgba(255,255,255,0.7) !important; display: flex !important; align-items: center; gap: 5px; }
.cw-status-dot { width: 6px; height: 6px; border-radius: 50%; background: #34D399 !important; animation: cwBlink 2s infinite; }
@keyframes cwBlink { 0%,100%{opacity:1} 50%{opacity:0.35} }
.cw-header-right { display: flex !important; gap: 6px; }
.cw-hbtn {
  width: 36px; height: 36px; border: none !important; border-radius: 10px;
  background: rgba(255,255,255,0.14) !important; color: #fff !important;
  display: flex !important; align-items: center; justify-content: center;
  cursor: pointer; font-size: 18px; transition: background 0.2s;
  -webkit-tap-highlight-color: transparent;
}
.cw-hbtn:active { background: rgba(255,255,255,0.28) !important; }
.cw-hbtn-close:active { background: rgba(239,68,68,0.45) !important; }

/* ═══ HISTORY DRAWER ═══ */
.cw-history {
  position: absolute; inset: 0; background: var(--cw-bg) !important;
  z-index: 10; display: flex !important; flex-direction: column;
}
.cw-history-top {
  display: flex !important; align-items: center; justify-content: space-between;
  padding: 16px; padding-top: calc(16px + env(safe-area-inset-top, 0px));
  border-bottom: 1px solid var(--cw-border) !important; background: var(--cw-card) !important;
}
.cw-history-top h4 { margin: 0; font-size: 15px; font-weight: 700; color: var(--cw-text) !important; }
.cw-history-top .cw-hbtn { background: #F3F4F6 !important; color: #6B7280 !important; }
.cw-history-list { overflow-y: auto; flex: 1; padding: 8px; -webkit-overflow-scrolling: touch; }
.cw-history-item {
  display: flex !important; align-items: center; justify-content: space-between;
  padding: 14px; border-radius: 12px; cursor: pointer;
  transition: background 0.15s; margin-bottom: 4px;
  -webkit-tap-highlight-color: transparent;
}
.cw-history-item:active { background: #F5F5F5 !important; }
.cw-history-item.active { background: var(--cw-primary-light) !important; border-left: 3px solid var(--cw-primary) !important; }
.cw-history-body { flex: 1; min-width: 0; }
.cw-history-title { display: block; font-size: 14px; font-weight: 500; color: var(--cw-text) !important;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.cw-history-time { display: block; font-size: 11px; color: #9CA3AF !important; margin-top: 3px; }
.cw-history-del {
  width: 32px; height: 32px; border: none !important; border-radius: 8px;
  background: transparent !important; color: #9CA3AF !important;
  cursor: pointer; display: flex !important; align-items: center; justify-content: center;
  font-size: 15px; flex-shrink: 0; -webkit-tap-highlight-color: transparent;
}
.cw-history-del:active { background: #FEE2E2 !important; color: #EF4444 !important; }
.cw-history-empty {
  display: flex !important; flex-direction: column; align-items: center;
  padding: 48px 20px; color: #9CA3AF !important; text-align: center;
}
.cw-history-empty i { font-size: 36px; margin-bottom: 12px; opacity: 0.4; }
.cw-history-empty p { margin: 0; font-size: 13px; }

/* ═══ MESSAGES ═══ */
.cw-messages {
  flex: 1; overflow-y: auto; padding: 16px;
  display: flex !important; flex-direction: column; gap: 14px;
  background: var(--cw-bg) !important;
  scroll-behavior: smooth; -webkit-overflow-scrolling: touch;
}
.cw-messages::-webkit-scrollbar { width: 4px; }
.cw-messages::-webkit-scrollbar-thumb { background: #ddd; border-radius: 4px; }

/* Welcome */
.cw-welcome {
  display: flex !important; flex-direction: column; align-items: center;
  text-align: center; padding: 28px 20px; gap: 10px;
}
.cw-welcome-icon {
  width: 60px; height: 60px; border-radius: 20px;
  background: var(--cw-gradient) !important;
  display: flex !important; align-items: center; justify-content: center;
  font-size: 30px; color: #fff !important;
  box-shadow: 0 8px 24px var(--cw-primary-glow);
}
.cw-welcome h3 { margin: 0; font-size: 20px; font-weight: 700; color: var(--cw-text) !important; }
.cw-welcome p { margin: 0; font-size: 14px; color: var(--cw-text-secondary) !important; line-height: 1.6; max-width: 300px; }
.cw-welcome-grid {
  display: grid !important; grid-template-columns: 1fr 1fr; gap: 8px;
  width: 100%; max-width: 320px; margin-top: 6px;
}

/* ═══ MESSAGE BUBBLES ═══ */
.cw-msg { display: flex !important; gap: 8px; animation: cwMsgIn 0.3s ease-out; }
@keyframes cwMsgIn { from { opacity:0; transform:translateY(8px); } to { opacity:1; transform:translateY(0); } }
.cw-msg-user { flex-direction: row-reverse; }
.cw-msg-ava {
  width: 30px; height: 30px; border-radius: 10px;
  background: var(--cw-gradient) !important;
  display: flex !important; align-items: center; justify-content: center;
  color: #fff !important; font-size: 15px; flex-shrink: 0; margin-top: 2px;
}
.cw-bubble {
  max-width: 88%; border-radius: var(--cw-radius); position: relative;
  padding: 12px 16px;
}
.cw-msg-user .cw-bubble {
  background: var(--cw-gradient) !important; color: #fff !important;
  border-bottom-right-radius: 4px;
}
.cw-msg-ai .cw-bubble {
  background: var(--cw-card) !important; color: var(--cw-text) !important;
  border-bottom-left-radius: 4px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.05), 0 1px 2px rgba(0,0,0,0.03);
  border: 1px solid var(--cw-border) !important;
}
.cw-time { display: block; font-size: 10px; opacity: 0.5; margin-top: 6px; text-align: right; color: inherit !important; }

/* ═══ PROSE — GPT Markdown Rendering ═══ */
.cw-prose { font-size: 14px; line-height: 1.7; word-break: break-word; color: inherit !important; }

.cw-prose :deep(.cw-h2) {
  font-size: 16px; font-weight: 700; margin: 10px 0 6px; color: var(--cw-primary) !important;
  padding-bottom: 4px; border-bottom: 2px solid var(--cw-primary-light);
}
.cw-prose :deep(.cw-h3) { font-size: 15px; font-weight: 700; margin: 8px 0 4px; color: var(--cw-text) !important; }
.cw-prose :deep(.cw-h4) { font-size: 14px; font-weight: 600; margin: 6px 0 3px; color: var(--cw-text) !important; }
.cw-prose :deep(.cw-h5) { font-size: 13px; font-weight: 600; margin: 4px 0 2px; color: var(--cw-text-secondary) !important; }

.cw-prose :deep(.cw-p) { margin: 0 0 2px; line-height: 1.7; }
.cw-prose :deep(.cw-spacer) { height: 6px; }

.cw-prose :deep(strong) { font-weight: 600; color: inherit !important; }
.cw-prose :deep(em) { font-style: italic; }

.cw-prose :deep(.cw-ul),
.cw-prose :deep(.cw-ol) {
  margin: 6px 0; padding-left: 0; list-style: none;
}
.cw-prose :deep(.cw-ul li),
.cw-prose :deep(.cw-ol li) {
  position: relative; padding-left: 18px; margin-bottom: 4px;
  font-size: 13.5px; line-height: 1.6;
}
.cw-prose :deep(.cw-ul li)::before {
  content: ''; position: absolute; left: 4px; top: 9px;
  width: 6px; height: 6px; border-radius: 50%;
  background: var(--cw-primary) !important;
}
.cw-prose :deep(.cw-ol) { counter-reset: cw-ol-counter; }
.cw-prose :deep(.cw-ol li) { counter-increment: cw-ol-counter; }
.cw-prose :deep(.cw-ol li)::before {
  content: counter(cw-ol-counter) '.'; position: absolute; left: 0; top: 0;
  font-weight: 700; font-size: 13px; color: var(--cw-primary) !important;
}

.cw-prose :deep(.cw-inline-code) {
  background: var(--cw-primary-light) !important; color: var(--cw-primary-dark) !important;
  padding: 2px 6px; border-radius: 4px;
  font-family: 'SF Mono', 'Fira Code', monospace; font-size: 12.5px;
}
.cw-prose :deep(.cw-code) {
  background: #1E293B !important; color: #E2E8F0 !important;
  padding: 12px 14px; border-radius: 10px; overflow-x: auto;
  font-size: 12px; margin: 8px 0; font-family: 'SF Mono', 'Fira Code', monospace;
}

.cw-prose :deep(.cw-bq) {
  border-left: 3px solid var(--cw-primary) !important;
  padding: 8px 12px; margin: 8px 0;
  background: var(--cw-primary-light) !important;
  border-radius: 0 10px 10px 0; font-size: 13.5px;
  color: var(--cw-text-secondary) !important;
}

.cw-prose :deep(.cw-hr) {
  border: none !important; border-top: 1px solid var(--cw-border) !important; margin: 10px 0;
}

.cw-prose :deep(.cw-link) {
  color: var(--cw-primary) !important; text-decoration: none;
  font-weight: 500; border-bottom: 1px dashed var(--cw-primary);
  transition: opacity 0.15s;
}
.cw-prose :deep(.cw-link:hover) { opacity: 0.75; }

.cw-prose :deep(.cw-table) {
  width: 100%; border-collapse: collapse; margin: 8px 0;
  font-size: 12.5px; border-radius: 8px; overflow: hidden;
}
.cw-prose :deep(.cw-table th) {
  background: var(--cw-primary-light) !important; color: var(--cw-primary-dark) !important;
  padding: 8px 10px; text-align: left; font-weight: 600; font-size: 12px;
  border-bottom: 2px solid #F0D5C0 !important;
}
.cw-prose :deep(.cw-table td) {
  padding: 7px 10px; border-bottom: 1px solid var(--cw-border) !important;
}
.cw-prose :deep(.cw-table tr:last-child td) { border-bottom: none !important; }

/* User bubble prose override */
.cw-msg-user .cw-prose :deep(.cw-ul li)::before { background: rgba(255,255,255,0.7) !important; }
.cw-msg-user .cw-prose :deep(.cw-ol li)::before { color: rgba(255,255,255,0.85) !important; }
.cw-msg-user .cw-prose :deep(.cw-inline-code) { background: rgba(255,255,255,0.2) !important; color: #fff !important; }

/* ═══ TYPING ═══ */
.cw-bubble-typing { display: flex !important; align-items: center; gap: 10px; padding: 14px 18px; }
.cw-dots { display: flex !important; gap: 4px; align-items: center; }
.cw-dots span {
  width: 7px; height: 7px; border-radius: 50%;
  background: var(--cw-primary) !important; opacity: 0.35;
  animation: cwBounce 1.4s infinite;
}
.cw-dots span:nth-child(2) { animation-delay: 0.15s; }
.cw-dots span:nth-child(3) { animation-delay: 0.3s; }
@keyframes cwBounce { 0%,60%,100%{transform:translateY(0);opacity:0.35} 30%{transform:translateY(-5px);opacity:1} }
.cw-typing-text { font-size: 12px; color: var(--cw-text-secondary) !important; font-style: italic; }

/* ═══ QUICK ACTIONS & CHIPS ═══ */
.cw-actions { display: flex !important; flex-wrap: wrap; gap: 6px; padding: 4px 0; }
.cw-chip {
  background: var(--cw-card) !important;
  border: 1.5px solid #E5E7EB !important;
  border-radius: 22px; padding: 9px 15px;
  font-size: 13px; color: var(--cw-text) !important;
  cursor: pointer; transition: all 0.2s;
  font-family: inherit; display: inline-flex !important;
  align-items: center; gap: 6px;
  -webkit-tap-highlight-color: transparent;
}
.cw-chip i { font-size: 15px; color: var(--cw-primary) !important; }
.cw-chip:active {
  background: var(--cw-primary-light) !important;
  border-color: var(--cw-primary) !important;
  color: var(--cw-primary-dark) !important;
}

/* ═══ INPUT ═══ */
.cw-input-area {
  padding: 12px 16px; padding-bottom: calc(12px + env(safe-area-inset-bottom, 0px));
  background: var(--cw-card) !important;
  border-top: 1px solid var(--cw-border) !important; flex-shrink: 0;
}
.cw-input-box {
  display: flex !important; align-items: flex-end; gap: 8px;
  background: #F3F4F6 !important; border-radius: 14px;
  padding: 6px 6px 6px 14px; transition: all 0.2s;
  border: 1.5px solid transparent !important;
}
.cw-input-box.focused {
  background: var(--cw-card) !important;
  border-color: var(--cw-primary) !important;
  box-shadow: 0 0 0 3px rgba(232,112,42,0.1);
}
.cw-textarea {
  flex: 1; border: none !important; background: transparent !important;
  font-size: 15px; color: var(--cw-text) !important;
  resize: none; outline: none !important; font-family: inherit;
  max-height: 100px; line-height: 1.5; padding: 7px 0;
  -webkit-appearance: none;
}
.cw-textarea::placeholder { color: #9CA3AF !important; }
.cw-send {
  width: 38px; height: 38px; border: none !important; border-radius: 11px;
  background: var(--cw-gradient) !important; color: #fff !important;
  display: flex !important; align-items: center; justify-content: center;
  font-size: 18px; cursor: pointer; flex-shrink: 0;
  transition: all 0.2s; -webkit-tap-highlight-color: transparent;
}
.cw-send:active:not(:disabled) { transform: scale(0.93); }
.cw-send:disabled { opacity: 0.3; cursor: not-allowed; }
.cw-disclaimer { text-align: center; margin: 6px 0 0; font-size: 10px; color: #B0B7C3 !important; }

/* ═══ TRANSITIONS ═══ */
.cw-slide-enter-active { transition: all 0.3s cubic-bezier(0.4,0,0.2,1); }
.cw-slide-leave-active { transition: all 0.2s cubic-bezier(0.4,0,0.2,1); }
.cw-slide-enter-from { opacity: 0; }
.cw-slide-leave-to { opacity: 0; }
.cw-slide-enter-from .cw-panel { transform: translateY(100%); }
.cw-slide-leave-to .cw-panel { transform: translateY(100%); }
@media (min-width: 640px) {
  .cw-slide-enter-from .cw-panel { transform: translateY(30px) scale(0.95); }
  .cw-slide-leave-to .cw-panel { transform: translateY(30px) scale(0.95); }
}
.cw-drawer-enter-active { transition: all 0.3s ease; }
.cw-drawer-leave-active { transition: all 0.2s ease; }
.cw-drawer-enter-from { transform: translateX(-100%); opacity: 0; }
.cw-drawer-leave-to { transform: translateX(-100%); opacity: 0; }
</style>
