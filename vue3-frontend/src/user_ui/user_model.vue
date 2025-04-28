<template>
    <div class="system-container">
      <!-- 左侧导航栏 -->
      <div class="sidebar">
        <div class="logo">AI System</div>
        <nav>
          <router-link to="/user/main" class="nav-item">
            <span class="icon">🏠</span>
            <span>首页</span>
          </router-link>
          <router-link to="/user/profile" class="nav-item">
            <span class="icon">👤</span>
            <span>个人信息</span>
          </router-link>
          <router-link to="/user/model" class="nav-item">
            <span class="icon">🧠</span>
            <span>模型管理</span>
          </router-link>
          <router-link to="/user/task" class="nav-item">
            <span class="icon">📋</span>
            <span>任务管理</span>
          </router-link>
  
          <!-- 退出登录按钮 -->
          <div class="logout-area">
            <button class="logout-button" @click="handleLogout">
              <span class="icon">🚪</span>
              <span>退出登录</span>
            </button>
          </div>
        </nav>
      </div>
  
      <!-- 右侧内容区 -->
      <div class="main-content">
        <!-- 欢迎信息 -->
        <div v-if="$route.path === '/user/main'" class="welcome-container">
          <h1>欢迎回来，{{ userInfo.username }}！</h1>
          <p>组织：{{ userInfo.organization }}</p>
          <p>当前系统信息：</p>
          <ul class="system-info">
            <li>🕒 登录时间：{{ loginTime }}</li>
            <li>📌 当前版本：v7.78d</li>
            <li>📊 待处理任务：{{ pendingTasks }} 个</li>
          </ul>
        </div>
  
        <!-- 填充内容 -->
        <div class="content-placeholder" v-else>
          <h2>模型管理</h2>
          <p>请选择左侧菜单中的功能模块以查看详细内容。</p>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import { useRouter } from 'vue-router'
  
  const router = useRouter()
  
  // 从 localStorage 获取用户信息
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  
  // 登录时间和待处理任务
  const loginTime = ref(new Date().toLocaleString())
  const pendingTasks = ref(3)
  
  // 退出登录处理
  const handleLogout = () => {
    // 清除用户状态
    localStorage.removeItem('authToken')
    localStorage.removeItem('userInfo')
  
    // 跳转到登录页并阻止返回
    router.replace('/').then(() => {
      window.location.reload() // 可选：完全重置应用状态
    })
  }
  </script>
  
  <style scoped>
  /* 页面布局 */
  .system-container {
    display: flex;
    min-height: 100vh;
    background-color: #f5f5f5;
  }
  
  /* 左侧导航栏 */
  .sidebar {
    width: 240px;
    background-color: #2c3e50;
    color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
    position: relative;
  }
  
  .logo {
    font-size: 1.5rem;
    font-weight: bold;
    text-align: center;
    margin-bottom: 20px;
  }
  
  .nav-item {
    display: flex;
    align-items: center;
    padding: 10px 15px;
    color: white;
    text-decoration: none;
    border-radius: 8px;
    margin-bottom: 10px;
    transition: background-color 0.3s ease;
  }
  
  .nav-item:hover {
    background-color: #34495e;
  }
  
  .nav-item .icon {
    margin-right: 10px;
  }
  
  .logout-area {
    position: absolute;
    bottom: 20px;
    width: calc(100% - 40px); /* 适配侧边栏padding */
  }
  
  .logout-button {
    width: 100%;
    padding: 15px 20px;
    background: none;
    border: none;
    color: #bdc3c7;
    text-align: left;
    cursor: pointer;
    display: flex;
    align-items: center;
    border-radius: 8px;
    transition: all 0.3s;
  }
  
  .logout-button:hover {
    background: #e74c3c;
    color: white;
  }
  
  /* 右侧内容区 */
  .main-content {
    flex: 1;
    padding: 20px;
    background-color: white;
    box-shadow: inset 0 0 10px rgba(0, 0, 0, 0.1);
  }
  
  /* 欢迎信息样式 */
  .welcome-container {
    background: white;
    padding: 30px;
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  }
  
  .welcome-container h1 {
    color: #2c3e50;
    margin-bottom: 20px;
  }
  
  .system-info {
    list-style: none;
    padding: 0;
  }
  
  .system-info li {
    padding: 10px 0;
    border-bottom: 1px solid #eee;
    display: flex;
    align-items: center;
    gap: 10px;
  }
  
  .system-info li:last-child {
    border-bottom: none;
  }
  
  /* 填充内容样式 */
  .content-placeholder {
    background: #f9f9f9;
    padding: 20px;
    border-radius: 10px;
    text-align: center;
    color: #7f8c8d;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  }
  
  .content-placeholder h2 {
    color: #2c3e50;
    margin-bottom: 10px;
  }
  </style>