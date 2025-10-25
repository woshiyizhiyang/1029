<template>
  <div class="character-animation">
    <div class="animation-container">
      <img
        v-for="(image, index) in images"
        :key="index"
        :src="image"
        :class="['mouth-image', { active: currentState === index }]"
        :alt="`口型状态${index}`"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useChatStore } from '@/store/chat'

const chatStore = useChatStore()

// 口型图片路径（需要放置5张图片到 public/images 目录）
const images = ref([
  '/images/mouth-0.jpg', // 闭口
  '/images/mouth-1.jpg', // 微张
  '/images/mouth-2.jpg', // 半张
  '/images/mouth-3.jpg', // 大张
  '/images/mouth-4.jpg'  // 完全张口
])

const currentState = ref(0)

// 监听store中的口型状态
watch(() => chatStore.currentMouthState, (newState) => {
  currentState.value = newState
})

// 预加载所有图片
const preloadImages = () => {
  images.value.forEach(src => {
    const img = new Image()
    img.src = src
  })
}

preloadImages()
</script>

<style scoped>
.character-animation {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.animation-container {
  position: relative;
  width: 100%;
  max-width: 720px;
  aspect-ratio: 720 / 1440;
  overflow: hidden;
}

.mouth-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
  opacity: 0;
  transition: opacity 0.05s ease-in-out;
  will-change: opacity;
}

.mouth-image.active {
  opacity: 1;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .character-animation {
    padding: 10px;
  }
}
</style>
