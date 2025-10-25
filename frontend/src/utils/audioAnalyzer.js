export class AudioAnalyzer {
  constructor() {
    this.audioContext = null
    this.analyser = null
    this.source = null
    this.audioElement = null
    this.animationFrameId = null
    this.onMouthStateChange = null
  }

  async init(audioBase64) {
    try {
      // 创建音频上下文
      this.audioContext = new (window.AudioContext || window.webkitAudioContext)()
      this.analyser = this.audioContext.createAnalyser()
      this.analyser.fftSize = 256
      
      // 解码Base64音频
      const audioData = this.base64ToArrayBuffer(audioBase64)
      const audioBuffer = await this.audioContext.decodeAudioData(audioData)
      
      // 创建音频元素
      this.audioElement = new Audio()
      this.audioElement.src = URL.createObjectURL(new Blob([audioData], { type: 'audio/mp3' }))
      
      // 创建音频源
      this.source = this.audioContext.createMediaElementSource(this.audioElement)
      this.source.connect(this.analyser)
      this.analyser.connect(this.audioContext.destination)
      
      return this.audioElement
    } catch (error) {
      console.error('Failed to initialize audio:', error)
      throw error
    }
  }

  play(onMouthStateChange) {
    if (!this.audioElement) return Promise.reject(new Error('Audio not initialized'))
    
    this.onMouthStateChange = onMouthStateChange
    this.startAnalyzing()
    
    return new Promise((resolve, reject) => {
      this.audioElement.onended = () => {
        this.stop()
        resolve()
      }
      
      this.audioElement.onerror = (error) => {
        this.stop()
        reject(error)
      }
      
      this.audioElement.play().catch(reject)
    })
  }

  startAnalyzing() {
    const bufferLength = this.analyser.frequencyBinCount
    const dataArray = new Uint8Array(bufferLength)
    
    const analyze = () => {
      if (!this.analyser) return
      
      this.animationFrameId = requestAnimationFrame(analyze)
      this.analyser.getByteFrequencyData(dataArray)
      
      // 计算平均音量
      const average = dataArray.reduce((sum, value) => sum + value, 0) / bufferLength
      
      // 映射到口型状态 (0-4)
      const mouthState = this.volumeToMouthState(average)
      
      if (this.onMouthStateChange) {
        this.onMouthStateChange(mouthState)
      }
    }
    
    analyze()
  }

  volumeToMouthState(volume) {
    // 将音量(0-255)映射到口型状态(0-4)
    // 音量范围映射规则:
    // 0-30: 闭口(0)
    // 30-80: 微张(1)
    // 80-130: 半张(2)
    // 130-180: 大张(3)
    // 180-255: 完全张口(4)
    
    if (volume < 30) return 0
    if (volume < 80) return 1
    if (volume < 130) return 2
    if (volume < 180) return 3
    return 4
  }

  stop() {
    if (this.animationFrameId) {
      cancelAnimationFrame(this.animationFrameId)
      this.animationFrameId = null
    }
    
    if (this.audioElement) {
      this.audioElement.pause()
      this.audioElement.currentTime = 0
    }
    
    if (this.source) {
      this.source.disconnect()
    }
    
    if (this.analyser) {
      this.analyser.disconnect()
    }
    
    if (this.audioContext && this.audioContext.state !== 'closed') {
      this.audioContext.close()
    }
    
    this.audioContext = null
    this.analyser = null
    this.source = null
    this.audioElement = null
  }

  base64ToArrayBuffer(base64) {
    const binaryString = window.atob(base64)
    const len = binaryString.length
    const bytes = new Uint8Array(len)
    for (let i = 0; i < len; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    return bytes.buffer
  }
}

export default AudioAnalyzer
