<template>
  <div class="app">
    <h1>📚 AI智能单词本</h1>
    
    <!-- 登录/注册 -->
    <div v-if="!token" class="auth">
      <div class="tabs">
        <button :class="{active: mode==='login'}" @click="mode='login'">登录</button>
        <button :class="{active: mode==='register'}" @click="mode='register'">注册</button>
      </div>
      <input v-model="username" placeholder="用户名">
      <input v-model="password" type="password" placeholder="密码">
      <button @click="mode==='login'?login():register()">{{mode==='login'?'登录':'注册'}}</button>
      <div class="error">{{authError}}</div>
    </div>

    <!-- 主界面 -->
    <div v-else class="main">
      <div class="header">
        <span>欢迎, {{user?.username}}</span>
        <button @click="logout">退出</button>
      </div>

      <!-- 查询 -->
      <div class="query">
        <h2>查询单词</h2>
        <div class="search">
          <input v-model="searchWord" @keyup.enter="query" placeholder="输入单词">
          <select v-model="provider">
            <option value="deepseek">DeepSeek</option>
            <option value="qwen">通义千问</option>
          </select>
          <button @click="query" :disabled="loading">{{loading?'查询中...':'查询'}}</button>
        </div>
        
        <div v-if="result" class="result">
          <h3>{{result.word}}</h3>
          <p>{{result.definition}}</p>
          <ul><li v-for="s in result.sentences">{{s}}</li></ul>
          <button @click="save" :disabled="saved">保存到单词本</button>
          <span class="source">来源: {{result.ai_provider}}</span>
        </div>
      </div>

      <!-- 单词列表 -->
      <div class="list">
        <h2>我的单词本</h2>
        <div class="pagination">
          <button @click="prevPage" :disabled="page===1">上一页</button>
          <span>{{page}}/{{totalPages}}</span>
          <button @click="nextPage" :disabled="page===totalPages">下一页</button>
        </div>
        <div v-for="w in words" class="word-item">
          <div><strong>{{w.word}}</strong> <button @click="del(w.id)">删除</button></div>
          <p>{{w.definition}}</p>
          <ul><li v-for="s in w.sentences">{{s}}</li></ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      token: localStorage.getItem('token') || '',
      user: null,
      mode: 'login',
      username: '',
      password: '',
      authError: '',
      searchWord: '',
      provider: 'deepseek',
      result: null,
      loading: false,
      saved: false,
      words: [],
      page: 1,
      pageSize: 10,
      total: 0
    }
  },
  computed: {
    totalPages() { return Math.ceil(this.total / this.pageSize) }
  },
  mounted() {
    if (this.token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
      this.loadUser()
      this.fetchWords()
    }
  },
  methods: {
    setAuth() {
      axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
    },
    loadUser() {
      const u = localStorage.getItem('user')
      if (u) this.user = JSON.parse(u)
    },
    async register() {
      try {
        const res = await axios.post('/api/register', {username: this.username, password: this.password})
        this.token = res.data.token
        this.user = res.data.user
        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.user))
        this.setAuth()
        this.authError = ''
        this.fetchWords()
      } catch(e) {
        this.authError = e.response?.data?.error || '注册失败'
      }
    },
    async login() {
      try {
        const res = await axios.post('/api/login', {username: this.username, password: this.password})
        this.token = res.data.token
        this.user = res.data.user
        localStorage.setItem('token', this.token)
        localStorage.setItem('user', JSON.stringify(this.user))
        this.setAuth()
        this.authError = ''
        this.fetchWords()
      } catch(e) {
        this.authError = e.response?.data?.error || '登录失败'
      }
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.clear()
      delete axios.defaults.headers.common['Authorization']
      this.result = null
      this.words = []
    },
    async query() {
      if (!this.searchWord) return
      this.loading = true
      try {
        const res = await axios.get('/api/word/query', {params: {word: this.searchWord, ai_provider: this.provider}})
        this.result = res.data
        this.saved = res.data.from_cache
      } catch(e) {
        alert('查询失败')
      } finally {
        this.loading = false
      }
    },
    async save() {
      try {
        await axios.post('/api/word/save', this.result)
        this.saved = true
        alert('保存成功')
        this.fetchWords()
      } catch(e) {
        alert(e.response?.data?.error || '保存失败')
      }
    },
    async fetchWords() {
      try {
        const res = await axios.get('/api/words', {params: {page: this.page, page_size: this.pageSize}})
        this.words = res.data.words
        this.total = res.data.total
      } catch(e) {
        console.error('获取失败')
      }
    },
    async del(id) {
      if (!confirm('确定删除？')) return
      await axios.delete(`/api/word/${id}`)
      this.fetchWords()
    },
    prevPage() { this.page--; this.fetchWords() },
    nextPage() { this.page++; this.fetchWords() }
  }
}
</script>

<style src="./style.css"></style>