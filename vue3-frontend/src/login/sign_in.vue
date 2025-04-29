<template>
  <div class="auth-page">
    <div class="auth-card">
      <!-- 标题 -->
      <div class="brand-header">
        <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="80" height="80" />
        <h1 class="brand-title">用户登录</h1>
      </div>

      <!-- 登录表单 -->
      <form class="custom-form" @submit.prevent="handleLogin">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.username" type="text" class="form-input" placeholder="请输入用户名" required />
        </div>

        <div class="form-group">
          <label>密码</label>
          <input v-model="form.password" type="password" class="form-input" placeholder="请输入密码" required />
        </div>
        <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
        <button type="submit" class="submit-btn login-btn">立即登录</button>

        <div class="switch-link">
          没有账号？
          <router-link to="/register">去注册</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
// 引入依赖
import axios from "axios"
import { ref } from "vue"
import { useRouter } from "vue-router"

// 初始化变量
const router = useRouter()
const form = ref({ username: "", password: "" })
const errorMessage = ref("")
const successMessage = ref("")
let messageTimer = null

// 清除消息
const clearAllMessages = () => {
  errorMessage.value = ""
  successMessage.value = ""
  if (messageTimer) {
    clearTimeout(messageTimer)
    messageTimer = null
  }
}

// 登录处理
const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    console.log("请输入用户名和密码")
    return
  }
  try {
    const response = await axios.post("http://localhost:8089/login", form.value)
    const result = response.data

    // 保存用户信息和令牌到 localStorage
    localStorage.setItem("authToken", result.token)
    localStorage.setItem("userInfo", JSON.stringify(result.user))

    successMessage.value = "登录成功，3秒后跳转..."
    messageTimer = setTimeout(() => {
      successMessage.value = ""
      router.push("/user/main") // 跳转到主页面
    }, 3000)
  } catch (error) {
    errorMessage.value = error.response?.data?.message || "用户名或密码错误"
    messageTimer = setTimeout(() => {
      errorMessage.value = ""
    }, 5000)
  }
}
</script>

<style scoped>
/* 页面布局 */
.auth-page {
  background: rgb(255, 255, 255);
  min-height: 100vh; /* 让页面高度占满整个视口 */
  display: flex; /* 使用 Flexbox 布局 */
  justify-content: center; /* 水平居中 */
  align-items: center; /* 垂直居中 */
  padding: 2rem; /* 添加内边距 */
}

/* 登录卡片 */
.auth-card {
  width: 100%;
  max-width: 420px; /* 设置最大宽度 */
  background: rgba(255, 255, 255, 0.95); /* 半透明背景 */
  border-radius: 16px; /* 圆角 */
  padding: 40px; /* 内边距 */
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1); /* 添加阴影 */
}

.brand-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo {
  margin: 0 auto 1rem;
  width: 100px; /* 调整 logo 大小 */
  height: 100px;
}

.brand-title {
  color: #2c3e50;
  margin-bottom: 1.5rem;
  font-size: 1.8rem;
}

/* 表单样式 */
.custom-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.rem;
}

.form-group label {
  color: #4a5568;
  font-size: 0.9rem;
  font-weight: 500;
}

.form-input,
.form-select {
  width: 90%;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #4caf50;
  box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.2);
}

.submit-btn {
  width: 100%;
  padding: 14px;
  background-color: #0a9814;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
}

.submit-btn:hover {
  background-color: #06b437;
  transform: translateY(-1px);
}

.switch-link {
  text-align: center;
  margin-top: 1.5rem;
  color: #666;
}

.switch-link a {
  color: #2196f3;
  text-decoration: none;
  font-weight: 500;
  margin-left: 0.5rem;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .auth-card {
    padding: 20px;
    max-width: 90%; /* 调整宽度以适应小屏幕 */
  }

  .brand-title {
    font-size: 1.5rem;
  }

  .submit-btn {
    font-size: 0.9rem;
    padding: 12px;
  }
}
</style>
