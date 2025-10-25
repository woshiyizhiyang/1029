// 示例口型图片说明文件
// 此文件用于说明如何准备动画图片

/*
口型图片要求:

1. 图片数量: 5张
2. 图片尺寸: 720×1440 (竖版, 宽高比 1:2)
3. 图片格式: JPG/PNG (推荐WebP格式优化体积)
4. 文件命名: 
   - mouth-0.jpg (或 mouth-0.png / mouth-0.webp)
   - mouth-1.jpg
   - mouth-2.jpg
   - mouth-3.jpg
   - mouth-4.jpg

5. 存放位置: frontend/public/images/

6. 各图片对应状态:
   - mouth-0: 闭口状态 (静音或极小声)
   - mouth-1: 微张状态 (轻声说话)
   - mouth-2: 半张状态 (正常音量)
   - mouth-3: 大张状态 (较大音量)
   - mouth-4: 完全张口 (最大音量)

7. 图片设计建议:
   - 保持角色形象一致,仅口型变化
   - 确保清晰度,避免模糊
   - 优化文件大小,建议每张<500KB
   - 使用压缩工具如TinyPNG优化

8. 如果使用WebP格式:
   - 需要提供JPG降级方案
   - 或在代码中只使用JPG格式

9. 测试准备:
   - 可以先使用占位图测试功能
   - 确保图片路径正确
   - 检查浏览器控制台无加载错误

示例文件结构:
frontend/
  public/
    images/
      mouth-0.jpg  <- 闭口
      mouth-1.jpg  <- 微张
      mouth-2.jpg  <- 半张
      mouth-3.jpg  <- 大张
      mouth-4.jpg  <- 完全张口
*/
