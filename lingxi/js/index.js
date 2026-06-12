// 配置marked解析规则，适配Prism高亮
marked.setOptions({
    highlight: function(code, lang) {
        if (lang && Prism.languages[lang]) {
            return Prism.highlight(code, Prism.languages[lang], lang);
        }
        return Prism.highlight(code, Prism.languages.markdown, 'markdown');
    },
    breaks: true, // 支持换行符
    gfm: true // 支持GitHub Flavored Markdown
});

// 核心配置：qwen-vl-plus视觉模型
const AI_CONFIG = {
    API_URL: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
    API_KEY: localStorage.getItem('LINGXI_API_KEY') || '',
    MODEL: "qwen-vl-plus",
    TEMPERATURE: 0.7
};

// DOM元素
const chatContent = document.getElementById('chat-content');
const messageInput = document.getElementById('message-input');
const sendBtn = document.getElementById('send-btn');
const chatHeader = document.getElementById('chat-header');
const suggestQuestions = document.getElementById('suggest-questions');
const clearBtn = document.getElementById('clear-btn');
const fileInput = document.getElementById('file-input');
const themeBtn = document.getElementById('theme-btn');
const uploadStopContainer = document.getElementById('upload-stop-container');
const inputImagePreview = document.getElementById('inputImagePreview');
const imageModal = document.getElementById('imageModal');
const modalImage = document.getElementById('modalImage');
const closeModal = document.getElementById('closeModal');

// 全局状态
let hasStartedChat = false;
let isLightMode = false;
let isGenerating = false;
let currentAbortController = null;
let currentAIElement = null;
let uploadedImages = [];
let sendDebounce = false; // 防止重复发送

// 优化滚动：AI回答时自动滚动到底部
function scrollToBottom(force = false) {
    // 强制滚动或AI生成时自动滚动
    if (force || isGenerating) {
        chatContent.scrollTop = chatContent.scrollHeight;
    }
}

function hideWelcome() {
    if (!hasStartedChat) {
        chatHeader.classList.add('hidden');
        suggestQuestions.classList.add('hidden');
        hasStartedChat = true;
    }
}

// 切换上传/停止按钮
function toggleUploadStop(showStop = false) {
    if (showStop) {
        uploadStopContainer.innerHTML = `
            <button class="stop-btn" id="stop-btn">⏹</button>
        `;
        setTimeout(() => {
            const stopBtn = document.getElementById('stop-btn');
            if (stopBtn) stopBtn.addEventListener('click', stopGenerate);
        }, 0);
    } else {
        uploadStopContainer.innerHTML = `
            <label class="func-btn" for="file-input">🖼️</label>
            <input type="file" id="file-input" accept="image/*">
        `;
        setTimeout(() => {
            const fileInput = document.getElementById('file-input');
            if (fileInput) fileInput.addEventListener('change', handleImageUpload);
        }, 0);
    }
}

// 停止生成
function stopGenerate() {
    if (isGenerating && currentAbortController) {
        currentAbortController.abort();
        isGenerating = false;
        sendDebounce = false;
        sendBtn.disabled = false;
        toggleUploadStop(false);
        if (currentAIElement) {
            currentAIElement.innerHTML += '<br><span style="color:var(--danger-color);font-style:italic;">（已停止）</span>';
            scrollToBottom(true);
        }
    }
}

// 格式化消息：Markdown转HTML + 代码块处理
function formatContent(content) {
    if (!content) return '';
    
    // 使用marked解析Markdown为HTML
    let htmlContent = marked.parse(content);
    
    // 为代码块添加复制按钮
    htmlContent = htmlContent.replace(/<pre><code class="language-([^"]+)">(.*?)<\/code><\/pre>/gs, (match, lang, code) => {
        // 解码HTML实体（&lt; → < 等）
        const decodedCode = decodeHTMLEntities(code);
        return `<div class="code-block">
                    <button class="copy-btn" data-code="${encodeURIComponent(decodedCode)}">复制</button>
                    <pre><code class="language-${lang}">${code}</code></pre>
                </div>`;
    });
    
    return htmlContent;
}

// HTML实体解码工具函数
function decodeHTMLEntities(text) {
    const textArea = document.createElement('textarea');
    textArea.innerHTML = text;
    return textArea.value;
}

// 初始化代码块
function initCodeBlocks() {
    Prism.highlightAll();
    document.querySelectorAll('.copy-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            const code = decodeURIComponent(this.dataset.code);
            navigator.clipboard.writeText(code).then(() => {
                this.textContent = '已复制';
                this.classList.add('copied');
                setTimeout(() => {
                    this.textContent = '复制';
                    this.classList.remove('copied');
                }, 2000);
            }).catch(err => {
                alert('复制失败：' + err.message);
            });
        });
    });
}

// 绑定图片预览
function bindImagePreview() {
    // 绑定聊天内容中的图片预览
    document.querySelectorAll('.image-preview').forEach(img => {
        img.addEventListener('click', () => {
            modalImage.src = img.src;
            imageModal.style.display = 'flex';
        });
    });
    
    // 新增：绑定输入框中的图片预览
    document.querySelectorAll('.input-image-item img').forEach(img => {
        img.addEventListener('click', (e) => {
            // 阻止事件冒泡，避免触发删除按钮的点击
            e.stopPropagation();
            modalImage.src = img.src;
            imageModal.style.display = 'flex';
        });
    });
    
    // 新增：确保删除按钮可以正常点击
    document.querySelectorAll('.remove-image').forEach(btn => {
        btn.addEventListener('click', (e) => {
            e.stopPropagation();
        });
    });
}

// 图文请求核心逻辑
async function callAI(prompt, images = []) {
    const apiKey = localStorage.getItem('LINGXI_API_KEY') || '';
    
    if (!apiKey || apiKey.trim() === '') {
        alert('请先在代码中配置有效的API Key！');
        return null;
    }

    const content = [
        { type: "text", text: prompt || "分析这张图片的内容" }
    ];

    for (const imgUrl of images) {
        if (imgUrl.startsWith('data:image/')) {
            const base64Data = imgUrl.split(',')[1];
            content.push({
                type: "image_url",
                image_url: {
                    url: `data:image/png;base64,${base64Data}`
                }
            });
        }
    }

    currentAbortController = new AbortController();
    try {
        const res = await fetch(AI_CONFIG.API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${apiKey}`
            },
            body: JSON.stringify({
                model: AI_CONFIG.MODEL,
                messages: [
                    { role: "user", content: content }
                ],
                stream: true,
                temperature: AI_CONFIG.TEMPERATURE,
                max_tokens: 2048
            }),
            signal: currentAbortController.signal
        });

        if (!res.ok) {
            const errText = await res.text();
            throw new Error(`API错误：${res.status} - ${errText}`);
        }
        return res.body;
    } catch (e) {
        if (e.name !== 'AbortError') {
            alert(`调用失败：${e.message}`);
        }
        return null;
    }
}

// 处理流式响应（优化自动滚动）
async function handleStream(stream, element) {
    if (!stream) return;
    const reader = stream.getReader();
    const decoder = new TextDecoder();
    let fullContent = '';
    try {
        while (true) {
            const { done, value } = await reader.read();
            if (done || !isGenerating) break;
            
            const chunk = decoder.decode(value);
            const lines = chunk.split('\n').filter(line => 
                line.trim() !== '' && line !== 'data: [DONE]'
            );
            
            for (const line of lines) {
                if (!line.startsWith('data: ')) continue;
                const data = line.substring(6);
                try {
                    const json = JSON.parse(data);
                    const delta = json.choices?.[0]?.delta?.content || '';
                    if (delta) {
                        fullContent += delta;
                        // 流式渲染时先简单处理换行，最终统一解析Markdown
                        element.innerHTML += delta.replace(/\n/g, '<br>');
                        // 每接收一点内容就滚动到底部
                        scrollToBottom(true);
                    }
                } catch (e) {
                    continue;
                }
            }
        }
    } finally {
        reader.releaseLock();
        isGenerating = false;
        sendDebounce = false;
        sendBtn.disabled = false;
        toggleUploadStop(false);
        // 最终渲染完整的Markdown内容
        element.innerHTML = formatContent(fullContent);
        initCodeBlocks();
        bindImagePreview();
        scrollToBottom(true);
    }
}

// 添加消息
function addMessage(text, isUser, imgs = []) {
    const msg = document.createElement('div');
    msg.className = `message ${isUser ? 'user' : 'ai'}`;
    const avatar = document.createElement('div');
    avatar.className = 'avatar';
    avatar.textContent = isUser ? '👤' : '🤖';
    const content = document.createElement('div');
    content.className = 'message-content';
    
    imgs.forEach(url => {
        const img = document.createElement('img');
        img.src = url;
        img.className = 'image-preview';
        content.appendChild(img);
        content.appendChild(document.createElement('br'));
    });
    
    if (isUser) {
        // 用户消息也解析Markdown
        content.innerHTML += formatContent(text);
        msg.append(avatar, content);
        chatContent.appendChild(msg);
        initCodeBlocks();
        bindImagePreview();
    } else {
        msg.append(avatar, content);
        chatContent.appendChild(msg);
        currentAIElement = content;
    }
    scrollToBottom(true);
}

// 发送消息
async function send() {
    // 防抖处理，防止重复发送
    if (sendDebounce) return;
    
    const text = messageInput.value.trim();
    if (!text && uploadedImages.length === 0) {
        alert('请输入问题，或上传图片后提问！');
        return;
    }
    
    // 设置防抖标记
    sendDebounce = true;
    sendBtn.disabled = true;
    hideWelcome();
    isGenerating = true;

    addMessage(text || '分析这张图片的内容', true, uploadedImages);
    messageInput.value = '';
    const imgs = [...uploadedImages];
    uploadedImages = [];
    inputImagePreview.innerHTML = '';

    const loading = document.createElement('div');
    loading.className = 'message ai';
    loading.innerHTML = `<div class="avatar">🤖</div><div class="message-content"><div class="loading"></div></div>`;
    chatContent.appendChild(loading);
    scrollToBottom(true);
    toggleUploadStop(true);

    try {
        chatContent.removeChild(loading);
        const stream = await callAI(text, imgs);
        if (stream) {
            addMessage('', false);
            await handleStream(stream, currentAIElement);
        } else {
            toggleUploadStop(false);
            isGenerating = false;
            sendDebounce = false;
            sendBtn.disabled = false;
        }
    } catch (e) {
        chatContent.removeChild(loading);
        addMessage(`错误：${e.message}`, false);
        toggleUploadStop(false);
        isGenerating = false;
        sendDebounce = false;
        sendBtn.disabled = false;
    }
}

// 快捷提问
function sendSuggestMessage(msg) {
    messageInput.value = msg;
    send();
}

// 清空对话（修改：添加滚动到顶部）
function clearChat() {
    stopGenerate();
    chatContent.innerHTML = '';
    hasStartedChat = false;
    chatHeader.classList.remove('hidden');
    suggestQuestions.classList.remove('hidden');
    messageInput.value = '';
    uploadedImages = [];
    inputImagePreview.innerHTML = '';
    // 滚动到页面最顶部
    window.scrollTo({
        top: 0,
        behavior: 'smooth' // 平滑滚动效果
    });
    scrollToBottom(true);
}

// 图片上传
function handleImageUpload(e) {
    const files = Array.from(e.target.files);
    files.forEach(file => {
        if (!file.type.startsWith('image/')) {
            alert('仅支持上传图片文件！');
            return;
        }
        const reader = new FileReader();
        reader.onload = ev => {
            const url = ev.target.result;
            uploadedImages.push(url);
            const item = document.createElement('div');
            item.className = 'input-image-item';
            item.innerHTML = `
                <img src="${url}" style="width:100%;height:100%;object-fit:cover;border-radius:8px;">
                <span class="remove-image" data-url="${url}">×</span>
            `;
            inputImagePreview.appendChild(item);
            
            // 为新添加的图片绑定预览事件
            const imgElement = item.querySelector('img');
            imgElement.addEventListener('click', (e) => {
                e.stopPropagation();
                modalImage.src = imgElement.src;
                imageModal.style.display = 'flex';
            });
            
            // 确保删除按钮可以正常点击，并且点击后滚动到顶部
            const removeBtn = item.querySelector('.remove-image');
            removeBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                uploadedImages = uploadedImages.filter(u => u !== url);
                item.remove();
                // 如果没有图片了，并且回到初始状态，滚动到顶部
                if (uploadedImages.length === 0 && !hasStartedChat) {
                    window.scrollTo({
                        top: 0,
                        behavior: 'smooth'
                    });
                }
            });
            
            hideWelcome();
        };
        reader.readAsDataURL(file);
    });
    e.target.value = '';
}

// 主题切换
function toggleTheme() {
    isLightMode = !isLightMode;
    document.documentElement.classList.toggle('light-mode', isLightMode);
    themeBtn.textContent = isLightMode ? '🌙' : '🔆';
    localStorage.setItem('theme', isLightMode ? 'light' : 'dark');
}

// 图片弹窗关闭
closeModal.onclick = () => imageModal.style.display = 'none';
imageModal.onclick = e => e.target === imageModal && (imageModal.style.display = 'none');
document.addEventListener('keydown', e => e.key === 'Escape' && (imageModal.style.display = 'none'));

// 绑定事件
sendBtn.addEventListener('click', send);
themeBtn.addEventListener('click', toggleTheme);
clearBtn.addEventListener('click', clearChat);
fileInput.addEventListener('change', handleImageUpload);
messageInput.addEventListener('keydown', e => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        send();
    }
});
messageInput.addEventListener('input', function() {
    this.style.height = 'auto';
    this.style.height = Math.min(this.scrollHeight, 120) + 'px';
});

// 初始化
window.addEventListener('load', () => {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'light') {
        isLightMode = true;
        document.documentElement.classList.add('light-mode');
        themeBtn.textContent = '🌙';
    }
    messageInput.focus();
    
    // 标题逐字渐变动画
    const titleSpans = document.querySelectorAll('.header-title span');
    titleSpans.forEach((span, index) => {
        span.style.animationDelay = `${index * 0.1}s`;
    });
    
    // 初始化滚动位置
    scrollToBottom(true);
});

// 窗口大小变化时重新计算高度并滚动
window.addEventListener('resize', () => {
    scrollToBottom(false);
});