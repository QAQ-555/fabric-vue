<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="brand-header">
        <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="80" height="80" />
        <h1 class="brand-title">用户注册</h1>
      </div>

      <form class="custom-form" @submit.prevent="handleRegister">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="form.username" type="text" class="form-input" placeholder="请输入用户名" required>
        </div>

        <div class="form-group">
          <label>密码</label>
          <input v-model="form.password" type="password" class="form-input" placeholder="请输入密码" required>
        </div>

        <div class="form-group">
          <label>选择组织</label>
          <select v-model="form.organization" class="form-select" required>
            <option value="" disabled>请选择组织</option>
            <option value="org1">组织一</option>
            <option value="org2">组织二</option>
          </select>
        </div>

        <button type="submit" class="submit-btn">
          注册
        </button>

        <div class="switch-link">
          已有账号？
          <router-link to="/">去登录</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = ref({
  username: '',
  password: '',
  organization: ''
})

const handleRegister = async () => {
  if (!form.value.username || !form.value.password|| !form.value.organization) {
    console.log('请填写完整信息')
    return
  }

  try {
    // 调用后端注册接口
    const response = await fetch('http://localhost:8089/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form.value)
    })

    const result = await response.json()

    if (response.ok) {
      // 注册成功，跳转到登录页面
      console.log('注册成功:', result.message)
      alert('注册成功，请登录！')
      router.push('/login')
    } else {
      // 注册失败，显示错误信息
      console.error('注册失败:', result.error)
      alert(`注册失败: ${result.error}`)
    }
  } catch (error) {
    console.error('注册请求失败:', error)
    alert('注册请求失败，请稍后重试！')
  }
}
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
  width: 100px;
  /* 调整 logo 大小 */
  height: 100px;
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
    max-width: 90%;
    /* 调整宽度以适应小屏幕 */
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