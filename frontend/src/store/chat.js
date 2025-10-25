import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useChatStore = defineStore('chat', () => {
  // 状态
  const messages = ref([])
  const sessionId = ref('')
  const isAISpeaking = ref(false)
  const isSending = ref(false)
  const currentMouthState = ref(0)
  const connectionState = ref('disconnected')
  const inputText = ref('')
  const showHistory = ref(false)

  // 计算属性
  const charCount = computed(() => inputText.value.length)
  const canSend = computed(() => {
    return !isSending.value && 
           inputText.value.trim() !== '' && 
           charCount.value <= 10000
  })
  const isCharLimitWarning = computed(() => charCount.value > 9000)
  const hasMessages = computed(() => messages.value.length > 0)

  // 方法
  function addMessage(message) {
    messages.value.push({
      ...message,
      timestamp: Date.now()
    })
  }

  function clearInput() {
    inputText.value = ''
  }

  function setSessionId(id) {
    sessionId.value = id
    // 保存到localStorage
    if (typeof window !== 'undefined') {
      localStorage.setItem('ai_chat_session_id', id)
    }
  }

  function loadSessionId() {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('ai_chat_session_id')
      if (saved) {
        sessionId.value = saved
      }
    }
  }

  function setConnectionState(state) {
    connectionState.value = state
  }

  function setSending(value) {
    isSending.value = value
  }

  function setAISpeaking(value) {
    isAISpeaking.value = value
  }

  function setMouthState(state) {
    currentMouthState.value = state
  }

  function toggleHistory() {
    showHistory.value = !showHistory.value
  }

  function resetMouth() {
    currentMouthState.value = 0
  }

  return {
    // 状态
    messages,
    sessionId,
    isAISpeaking,
    isSending,
    currentMouthState,
    connectionState,
    inputText,
    showHistory,
    // 计算属性
    charCount,
    canSend,
    isCharLimitWarning,
    hasMessages,
    // 方法
    addMessage,
    clearInput,
    setSessionId,
    loadSessionId,
    setConnectionState,
    setSending,
    setAISpeaking,
    setMouthState,
    toggleHistory,
    resetMouth
  }
})
