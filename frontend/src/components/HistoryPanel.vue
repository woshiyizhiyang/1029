<template>
  <div class="history-panel" v-if="show">
    <div class="panel-header">
      <h3>对话历史</h3>
      <button class="close-button" @click="close">
        <i class="fas fa-times"></i>
      </button>
    </div>
    <div class="panel-content">
      <div
        v-for="(item, index) in history"
        :key="index"
        :class="['history-item', item.role]"
      >
        <div class="history-bubble">
          <div class="history-content">{{ item.content }}</div>
          <div class="history-time">{{ formatTime(item.created_at) }}</div>
        </div>
      </div>
      <div v-if="history.length === 0" class="empty-state">
        暂无对话记录
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  history: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close'])

const close = () => {
  emit('close')
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}
</script>

<style scoped>
.history-panel {
  position: fixed;
  top: 0;
  right: 0;
  width: 100%;
  max-width: 400px;
  height: 100%;
  background-color: #ffffff;
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #e0e0e0;
  background-color: #2196F3;
  color: white;
}

.panel-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.close-button {
  background: none;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  padding: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
}

.history-item {
  display: flex;
  margin-bottom: 15px;
}

.history-item.user {
  justify-content: flex-end;
}

.history-item.assistant,
.history-item.system {
  justify-content: flex-start;
}

.history-bubble {
  max-width: 80%;
  padding: 10px 14px;
  border-radius: 10px;
}

.history-item.user .history-bubble {
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
}

.history-item.assistant .history-bubble,
.history-item.system .history-bubble {
  background-color: #BBDEFB;
}

.history-content {
  font-size: 14px;
  line-height: 1.4;
  color: #333;
  word-wrap: break-word;
}

.history-time {
  font-size: 10px;
  color: #999;
  margin-top: 4px;
  text-align: right;
}

.empty-state {
  text-align: center;
  color: #999;
  padding: 40px 20px;
  font-size: 14px;
}

@media (max-width: 768px) {
  .history-panel {
    max-width: 100%;
  }
}
</style>
