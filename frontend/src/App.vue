<template>
  <div class="chat-container">
    <!-- 顶部通知 -->
    <div v-if="notification" :class="['notification', notification.type]">
      {{ notification.message }}
    </div>

    <!-- 历史按钮 -->
    <button
      v-if="chatStore.hasMessages"
      class="history-button"
      @click="toggleHistory"
    >
      <i class="fas fa-history"></i>
    </button>

    <!-- 动画形象区域 -->
    <div class="animation-area">
      <CharacterAnimation />
    </div>

    <!-- 对话区域 -->
    <div class="chat-area">
      <MessageList />
      <InputBox @send="handleSend" @stop="handleStop" />
    </div>

    <!-- 历史面板 -->
    <HistoryPanel
      :show="chatStore.showHistory"
      :history="historyData"
      @close="toggleHistory"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useChatStore } from '@/store/chat'
import CharacterAnimation from '@/components/CharacterAnimation.vue'
import MessageList from '@/components/MessageList.vue'
import InputBox from '@/components/InputBox.vue'
import HistoryPanel from '@/components/HistoryPanel.vue'
import WebSocketClient from '@/api/websocket'
import AudioAnalyzer from '@/utils/audioAnalyzer'

const chatStore = useChatStore()
const wsClient = ref(null)
const audioAnalyzer = ref(null)
const notification = ref(null)
const historyData = ref([])
const currentAIMessage = ref('')

onMounted(async () => {
  // 加载会话ID
  chatStore.loadSessionId()
  
  // 连接WebSocket
  await connectWebSocket()
  
  // 请求欢迎语
  requestWelcome()
})

onUnmounted(() => {
  if (wsClient.value) {
    wsClient.value.close()
  }
  if (audioAnalyzer.value) {
    audioAnalyzer.value.stop()
  }
})

// 连接WebSocket
const connectWebSocket = async () => {
  try {
    wsClient.value = new WebSocketClient()
    
    // 注册消息处理器
    wsClient.value.on('welcome', handleWelcome)
    wsClient.value.on('ai_text_chunk', handleTextChunk)
    wsClient.value.on('ai_text_complete', handleTextComplete)
    wsClient.value.on('audio_data', handleAudioData)
    wsClient.value.on('error', handleError)
    wsClient.value.on('stop_ack', handleStopAck)
    wsClient.value.on('history_data', handleHistoryData)
    
    await wsClient.value.connect('ws://localhost:8080/ws')
    chatStore.setConnectionState('connected')
    showNotification('已连接到服务器', 'info')
  } catch (error) {
    console.error('WebSocket connection failed:', error)
    chatStore.setConnectionState('disconnected')
    showNotification('连接服务器失败', 'error')
  }
}

// 请求欢迎语
const requestWelcome = () => {
  if (!wsClient.value) return
  
  const sessionId = chatStore.sessionId || generateSessionId()
  chatStore.setSessionId(sessionId)
  
  chatStore.setSending(true)
  wsClient.value.send({
    type: 'get_welcome',
    session_id: sessionId
  })
}

// 处理欢迎语
const handleWelcome = async (message) => {
  chatStore.addMessage({
    role: 'system',
    content: message.content
  })
  
  // 播放欢迎语音
  if (message.audioBase64) {
    await playAudio(message.audioBase64)
  }
  
  chatStore.setSending(false)
}

// 处理文本片段
const handleTextChunk = (message) => {
  currentAIMessage.value += message.content
  
  // 更新最后一条AI消息
  const messages = chatStore.messages
  if (messages.length > 0 && messages[messages.length - 1].role === 'assistant') {
    messages[messages.length - 1].content = currentAIMessage.value
  } else {
    chatStore.addMessage({
      role: 'assistant',
      content: currentAIMessage.value
    })
  }
}

// 处理文本完成
const handleTextComplete = (message) => {
  currentAIMessage.value = ''
}

// 处理音频数据
const handleAudioData = async (message) => {
  if (message.audioBase64) {
    await playAudio(message.audioBase64)
  }
  chatStore.setSending(false)
}

// 播放音频
const playAudio = async (audioBase64) => {
  try {
    chatStore.setAISpeaking(true)
    
    // 创建音频分析器
    if (audioAnalyzer.value) {
      audioAnalyzer.value.stop()
    }
    audioAnalyzer.value = new AudioAnalyzer()
    
    await audioAnalyzer.value.init(audioBase64)
    await audioAnalyzer.value.play((mouthState) => {
      chatStore.setMouthState(mouthState)
    })
    
    // 播放完成，重置口型
    chatStore.resetMouth()
    chatStore.setAISpeaking(false)
  } catch (error) {
    console.error('Audio playback failed:', error)
    showNotification('音频播放失败', 'warning')
    chatStore.setAISpeaking(false)
    chatStore.resetMouth()
  }
}

// 处理错误
const handleError = (message) => {
  showNotification(message.errorMsg, 'error')
  chatStore.setSending(false)
  chatStore.setAISpeaking(false)
  chatStore.resetMouth()
}

// 处理停止确认
const handleStopAck = () => {
  if (audioAnalyzer.value) {
    audioAnalyzer.value.stop()
  }
  chatStore.setSending(false)
  chatStore.setAISpeaking(false)
  chatStore.resetMouth()
  currentAIMessage.value = ''
}

// 处理历史数据
const handleHistoryData = (message) => {
  try {
    historyData.value = JSON.parse(message.content)
  } catch (error) {
    console.error('Failed to parse history data:', error)
  }
}

// 发送消息
const handleSend = () => {
  const text = chatStore.inputText.trim()
  if (!text || !wsClient.value) return
  
  // 添加用户消息
  chatStore.addMessage({
    role: 'user',
    content: text
  })
  
  // 发送到服务器
  wsClient.value.send({
    type: 'user_message',
    session_id: chatStore.sessionId,
    content: text
  })
  
  // 清空输入框
  chatStore.clearInput()
  chatStore.setSending(true)
  currentAIMessage.value = ''
}

// 停止
const handleStop = () => {
  if (!wsClient.value) return
  
  wsClient.value.send({
    type: 'stop',
    session_id: chatStore.sessionId
  })
  
  // 立即停止音频
  if (audioAnalyzer.value) {
    audioAnalyzer.value.stop()
  }
  chatStore.resetMouth()
}

// 切换历史面板
const toggleHistory = () => {
  if (!chatStore.showHistory) {
    // 请求历史数据
    wsClient.value.send({
      type: 'get_history',
      session_id: chatStore.sessionId
    })
  }
  chatStore.toggleHistory()
}

// 显示通知
const showNotification = (message, type = 'info') => {
  notification.value = { message, type }
  setTimeout(() => {
    notification.value = null
  }, 3000)
}

// 生成会话ID
const generateSessionId = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}
</script>

<style scoped>
.chat-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  position: relative;
}

.notification {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  z-index: 2000;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  animation: slideDown 0.3s ease-out;
}

.notification.info {
  background-color: #2196F3;
  color: white;
}

.notification.warning {
  background-color: #FF9800;
  color: white;
}

.notification.error {
  background-color: #F44336;
  color: white;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.history-button {
  position: fixed;
  top: 20px;
  right: 20px;
  width: 48px;
  height: 48px;
  border: none;
  border-radius: 50%;
  background-color: #2196F3;
  color: white;
  font-size: 20px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.history-button:hover {
  background-color: #1976D2;
  transform: scale(1.05);
}

.animation-area {
  height: 60vh;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chat-area {
  height: 40vh;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  border-top: 2px solid #e0e0e0;
}

@media (max-width: 768px) {
  .animation-area {
    height: 55vh;
  }
  
  .chat-area {
    height: 45vh;
  }
}
</style>
