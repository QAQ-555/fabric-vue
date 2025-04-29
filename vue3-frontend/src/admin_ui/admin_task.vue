<template>
  <div class="system-container">
    <!-- 左侧导航栏 -->
    <div class="sidebar">
      <div class="logo">AI System</div>
      <nav>
        <router-link to="/admin/main" class="nav-item">
          <span class="icon">🏠</span>
          <span>首页</span>
        </router-link>
        <router-link to="/admin/profile" class="nav-item">
          <span class="icon">👤</span>
          <span>个人信息</span>
        </router-link>
        <router-link to="/admin/model" class="nav-item">
          <span class="icon">🧠</span>
          <span>模型管理</span>
        </router-link>
        <router-link to="/admin/task" class="nav-item">
          <span class="icon">📋</span>
          <span>任务管理</span>
        </router-link>
        <router-link to="/admin/user" class="nav-item">
          <span class="icon">👥</span>
          <span>用户管理</span>
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

      <!-- 任务管理 -->
      <div class="content-placeholder" v-else>
        <h2>任务管理</h2>
        <div v-if="tasks.length > 0" class="tasks-container">
          <table class="tasks-table">
            <thead>
              <tr>
                <th>#</th>
                <th>任务ID</th>
                <th>根模型cid</th>
                <th>奖励</th>
                <th>发布者</th>
                <th>是否完成</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(task, index) in tasks" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ task.ID }}</td>
                <td>{{ task.rootModelHash }}</td>
                <td>{{ task.bonus }}</td>
                <td>{{ task.postedUser }}</td>
                <td>{{ task.isComplete ? '是' : '否' }}</td>
                <td>
                  <button @click="acceptTask(task.ID)">接受任务</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else>
          <p>暂无任务。</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 引入依赖
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

// 初始化变量
const router = useRouter()
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
const loginTime = ref(new Date().toLocaleString())
const pendingTasks = ref(3)
const tasks = ref([])

// 获取任务列表
const getAllTasks = async () => {
  try {
    const response = await axios.post("http://localhost:8089/get_all_task");
    tasks.value = response.data.tasks || [];
    console.log("获取的任务列表:", response.data);
  } catch (error) {
    console.error("获取任务失败:", error);
    alert("获取任务失败，请稍后重试！");
  }
}

// 接受任务
const acceptTask = async (taskId) => {
  const task = tasks.value.find(t => t.ID === taskId);
  if (task && task.isComplete) {
    alert(`任务 ${taskId} 已完成，无法接受！`);
    return;
  }

  try {
    const response = await axios.post("http://localhost:8089/accept_task", {
      taskId,
      username: userInfo.value.username
    });
    alert(`任务 ${taskId} 接受成功！`);
    // 可选：刷新任务列表
    getAllTasks();
  } catch (error) {
    console.error("接受任务失败:", error);
    alert("接受任务失败，请稍后重试！");
  }
}

// 页面加载时获取任务列表
onMounted(() => {
  getAllTasks()
})

// 退出登录
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

/* 任务表格样式 */
.tasks-container {
  margin-top: 20px;
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  max-height: 400px; /* 固定高度 */
  overflow-y: auto; /* 添加滚动条 */
}

.tasks-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

.tasks-table th,
.tasks-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: center;
}

.tasks-table th {
  background-color: #f4f4f4;
  color: #333;
  font-weight: bold;
}

.tasks-table tr:nth-child(even) {
  background-color: #f9f9f9;
}

.tasks-table tr:hover {
  background-color: #f1f1f1;
}

.tasks-table button {
  padding: 5px 10px;
  border: none;
  background-color: #4caf50;
  color: white;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.tasks-table button:hover {
  background-color: #45a049;
}
</style>