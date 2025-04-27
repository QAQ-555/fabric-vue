<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="brand-header">
        <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="80" height="80" />
        <h1 class="brand-title">用户登录</h1>
      </div>

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
import axios from "axios";
import { ref } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();
const form = ref({
  username: "",
  password: "",
});
const errorMessage = ref("");
const successMessage = ref("");
let messageTimer = null;

const clearAllMessages = () => {
  errorMessage.value = "";
  successMessage.value = "";
  if (messageTimer) {
    clearTimeout(messageTimer);
    messageTimer = null;
  }
};

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    console.log("请输入用户名和密码");
    return;
  }
  try {
    const response = await axios.post(
      "http://localhost:8089/login",
      form.value
    );
    successMessage.value = "登录成功，3秒后跳转...";
    messageTimer = setTimeout(() => {
      successMessage.value = "";
      router.push("/dashboard"); // 跳转操作
    }, 3000);
  } catch (error) {
    errorMessage.value = error.response?.data?.message || "用户名或密码错误";
    messageTimer = setTimeout(() => {
      errorMessage.value = "";
    }, 5000); // 错误提示5秒后消失
  }
};
</script>

<style scoped>
.auth-page {
  background: rgb(255, 255, 255);
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
}

.auth-card {
  width: 100%;
  max-width: 420px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.brand-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo {
  margin: 0 auto 1rem;
}

.brand-title {
  color: #2c3e50;
  margin-bottom: 1.5rem;
  font-size: 1.8rem;
}

.custom-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  color: #4a5568;
  font-size: 0.9rem;
  font-weight: 500;
}

.form-input,
.form-select {
  width: 100%;
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

.form-select {
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 1rem center;
  background-size: 1em;
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

.success-message {
  color: green;
  margin: 10px 0;
  padding: 10px;
  background: #e6ffe6;
  border-radius: 4px;
}

.error-message {
  color: red;
  font-size: 12px;
  margin-top: 5px;
  height: 20px;
}
</style>
