<template>
  <div class="input-box">
    <div class="input-wrapper">
      <textarea
        v-model="chatStore.inputText"
        :disabled="chatStore.isSending"
        :placeholder="placeholder"
        class="input-textarea"
        @keydown.enter.exact="handleEnter"
        rows="1"
      ></textarea>
      <button
        :class="['send-button', { sending: chatStore.isSending }]"
        :disabled="!chatStore.canSend && !chatStore.isSending"
        @click="handleButtonClick"
      >
        <i :class="buttonIcon"></i>
      </button>
    </div>
    <div class="char-counter" :class="{ warning: chatStore.isCharLimitWarning }">
      {{ chatStore.charCount }}/10000
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useChatStore } from '@/store/chat'

const chatStore = useChatStore()

const emit = defineEmits(['send', 'stop'])

const placeholder = computed(() => {
  if (chatStore.isSending) return '正在处理中...'
  return '请输入你的问题...'
})

const buttonIcon = computed(() => {
  return chatStore.isSending ? 'fas fa-stop' : 'fas fa-paper-plane'
})

const handleEnter = (e) => {
  if (!e.shiftKey) {
    e.preventDefault()
    if (chatStore.canSend) {
      emit('send')
    }
  }
}

const handleButtonClick = () => {
  if (chatStore.isSending) {
    emit('stop')
  } else if (chatStore.canSend) {
    emit('send')
  }
}
</script>

<style scoped>
.input-box {
  padding: 15px;
  background-color: #ffffff;
  border-top: 1px solid #e0e0e0;
}

.input-wrapper {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.input-textarea {
  flex: 1;
  min-height: 40px;
  max-height: 120px;
  padding: 10px 15px;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  font-size: 15px;
  font-family: inherit;
  resize: none;
  outline: none;
  transition: border-color 0.3s;
}

.input-textarea:focus {
  border-color: #2196F3;
}

.input-textarea:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.send-button {
  width: 44px;
  height: 44px;
  border: none;
  border-radius: 50%;
  background-color: #2196F3;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  flex-shrink: 0;
}

.send-button:hover:not(:disabled) {
  background-color: #1976D2;
  transform: scale(1.05);
}

.send-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.send-button.sending {
  background-color: #F44336;
}

.send-button.sending:hover {
  background-color: #D32F2F;
}

.send-button i {
  font-size: 16px;
}

.char-counter {
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.char-counter.warning {
  color: #F44336;
  font-weight: bold;
}
</style>
