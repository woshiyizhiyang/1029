<template>
  <div class="message-list" ref="messageContainer">
    <div
      v-for="(message, index) in messages"
      :key="index"
      :class="['message-item', message.role]"
    >
      <div class="message-bubble">
        <div class="message-content">{{ message.content }}</div>
        <div class="message-time">{{ formatTime(message.timestamp) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { useChatStore } from '@/store/chat'

const chatStore = useChatStore()
const messageContainer = ref(null)

const messages = ref([])

// 监听store中的消息变化
watch(() => chatStore.messages, (newMessages) => {
  messages.value = newMessages
  scrollToBottom()
}, { deep: true })

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}
</script>

<style scoped>
.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-item {
  display: flex;
  width: 100%;
}

.message-item.user {
  justify-content: flex-end;
}

.message-item.assistant,
.message-item.system {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  word-wrap: break-word;
}

.message-item.user .message-bubble {
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.message-item.assistant .message-bubble,
.message-item.system .message-bubble {
  background-color: #BBDEFB;
}

.message-content {
  font-size: 15px;
  line-height: 1.5;
  color: #333;
}

.message-time {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
  text-align: right;
}

/* 滚动条样式 */
.message-list::-webkit-scrollbar {
  width: 4px;
}

.message-list::-webkit-scrollbar-track {
  background: transparent;
}

.message-list::-webkit-scrollbar-thumb {
  background: #ccc;
  border-radius: 2px;
}
</style>
