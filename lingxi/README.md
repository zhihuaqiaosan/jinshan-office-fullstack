## 灵犀AI对话助手

### 一、项目基本信息

**姓名：陈慧**

**学校：华中师范大学**

**学号：2023214479**

**时间：2026.3.25**

### 二、已实现功能

| 任务模块                  | 完成状态 | 核心说明                                                     |
| ------------------------- | -------- | :----------------------------------------------------------- |
| 首页 UI 展示              | ✅ 完成   | 渐变色标题、4 张彩色快捷卡片、点击卡片自动发送提问           |
| AI 对话功能（阿里云百炼） | ✅ 完成   | 接入通义千问 qwen-vl-plus 模型，支持纯文本 / 图文混合对话，流式响应打字机效果 |
| Markdown 渲染             | ✅ 完成   | 支持标题、列表、表格、引用、代码块等全量 Markdown 语法解析   |
| 主题切换                  | ✅ 完成   | 深色 / 浅色模式切换，localStorage 持久化保存，刷新页面不丢失 |
| 底部输入区                | ✅ 完成   | 固定底部、图片上传预览、停止生成、清除对话、Enter 发送等功能 |
| 代码块增强                | ✅ 完成   | 语法高亮、右上角复制按钮、复制成功提示                       |
| 响应式适配                | ✅ 完成   | 适配 PC / 平板 / 手机等不同尺寸设备，卡片布局自动调整        |

### 三、核心技术实现

#### 1.阿里云百炼API接入

选择`qwen-vl-plus`多模态模型，同时支持文本和图片输入

请求参数设置`stream: true`开启流式响应，实现打字机效果

API Key 本地存储，避免硬编码

使用`AbortController`实现生成过程的中断（停止按钮功能）

#### 2.流式响应与打字机效果

```js
// 核心逻辑：逐块读取流式数据并实时渲染
async function handleStream(stream, element) {
    const reader = stream.getReader();
    const decoder = new TextDecoder();
    let fullContent = '';
    while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        const chunk = decoder.decode(value);
        // 解析每一块增量内容并追加到页面
        const delta = JSON.parse(chunk.substring(6)).choices[0]?.delta?.content || '';
        if (delta) {
            fullContent += delta;
            element.innerHTML += delta.replace(/\n/g, '<br>'); // 实时渲染
            scrollToBottom(true); // 自动滚动到底部
        }
    }
}
```

#### 3.Markdown渲染与代码高亮

使用`marked.js`解析 Markdown 为 HTML，配置自定义高亮规则

集成`Prism.js`实现代码语法高亮，支持 JavaScript/CSS/HTML/Python 等主流语言

为代码块添加自定义复制按钮，通过`navigator.clipboard`实现一键复制

#### 4.主题切换持久化

```js
// 初始化时读取本地存储的主题
window.addEventListener('load', () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'light') {
        document.documentElement.classList.add('light-mode');
        themeBtn.textContent = '🌙';
    }
});
// 切换主题并保存
function toggleTheme() {
    isLightMode = !isLightMode;
    document.documentElement.classList.toggle('light-mode', isLightMode);
    localStorage.setItem('theme', isLightMode ? 'light' : 'dark');
}
```

#### 5.图片上传与预览

通过`FileReader`将图片转为 base64 格式，支持预览和删除

点击预览图弹出大图查看弹窗，支持 ESC 键关闭

多模态请求时将 base64 图片数据传入 API，实现图文对话

### 四、遇到的问题与解决思路

#### 1.流式响应乱码

**问题**：首次解析流式数据时出现 JSON 解析错误或乱码

**解决**：过滤无效数据行（`data: [DONE]`），使用`TextDecoder`正确解码二进制流，增加异常捕获

#### 2.图片上传后API请求失效

**问题**：base64 图片数据过大导致请求超时

**解决**：优化图片处理逻辑，确保多模态请求格式符合阿里云百炼规范，增加请求超时处理

#### 3.主题切换样式不一致

**问题**：部分元素在浅色模式下样式未适配

**解决**：使用 CSS 变量统一管理主题色，为所有核心样式添加`light-mode`前缀的适配规则