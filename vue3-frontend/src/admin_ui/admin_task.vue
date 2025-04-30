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
                <th>轮数</th>
                <th>下一轮模型id</th>
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
                <td>{{ task.round }}</td>
                <td>{{ task.nextRoundTaskID }}</td>
                <td>{{ task.isComplete ? '是' : '否' }}</td>
                <td>
                  <button @click="acceptTask(task.ID)">接受任务</button>
                  <button @click="deleteTask(task.ID)">删除任务</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else>
          <p>暂无任务。</p>
        </div>

        <!-- 新增模块：管理由自己发布的任务 -->
        <div class="my-tasks-container" v-if="myTasks.length > 0">
          <h3>我发布的任务</h3>
          <table class="tasks-table">
            <thead>
              <tr>
                <th>#</th>
                <th>任务ID</th>
                <th>根模型cid</th>
                <th>奖励</th>
                <th>轮数</th>
                <th>下一轮模型id</th>
                <th>是否完成</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(task, index) in myTasks" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ task.ID }}</td>
                <td>{{ task.rootModelHash }}</td>
                <td>{{ task.bonus }}</td>
                <td>{{ task.round }}</td>
                <td>{{ task.nextRoundTaskID }}</td>
                <td>{{ task.isComplete ? '是' : '否' }}</td>
                <td>
                  <button @click="startNextRound(task.ID)">进入下一轮</button>
                  <button @click="finishTask(task.ID)">完成任务</button>
                </td>
              </tr>
            </tbody>
          </table>
          <button class="new-task-button" @click="showNewTaskModal = true">发布新任务</button>
        </div>
        <div v-else>
          <p>暂无我发布的任务。</p>
          <button class="new-task-button" @click="showNewTaskModal = true">发布新任务</button>
        </div>
      </div>
    </div>
  </div>

  <!-- 新增弹窗：发布新任务 -->
  <div v-if="showNewTaskModal" class="modal-overlay">
    <div class="modal-content">
      <h3>发布新任务</h3>
      <form @submit.prevent="submitNewTask">
        <div class="form-group">
          <label for="bonus">代币奖励值：</label>
          <input type="number" id="bonus" v-model="newTaskBonus" required />
        </div>
        <div class="form-group">
          <label for="rootModelId">根模型 ID：</label>
          <input type="text" id="rootModelId" v-model="newTaskRootModelId" required />
        </div>
        <div class="form-actions">
          <button type="submit" class="submit-button">提交</button>
          <button type="button" class="cancel-button" @click="showNewTaskModal = false">取消</button>
        </div>
      </form>
    </div>
  </div>

  <!-- 修改弹窗：进入下一轮 -->
  <div v-if="showNextRoundModal" class="modal-overlay">
    <div class="modal-content">
      <h3>进入下一轮</h3>
      <form @submit.prevent="submitNextRound">
        <div class="form-group">
          <label for="nextRoundRootModelId">根模型 ID：</label>
          <input type="text" id="nextRoundRootModelId" v-model="nextRoundRootModelId" required />
        </div>
        <div class="form-actions">
          <button type="submit" class="submit-button">提交</button>
          <button type="button" class="cancel-button" @click="cancelNextRound">取消</button>
        </div>
      </form>
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
const myTasks = ref([]);

// 获取任务列表
const getAllTasks = async () => {
  try {
    const response = await axios.post("http://localhost:8089/get_all_task");
    tasks.value = response.data.tasks || [];
    // 筛选出由当前用户发布的任务
    myTasks.value = tasks.value.filter(task => task.postedUser === userInfo.value.username);
    console.log("获取的任务列表:", tasks.value);
    console.log("我发布的任务:", myTasks.value);
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
// 删除
const deleteTask = async (taskId) => {
  const task = tasks.value.find(t => t.ID === taskId);
  try {
    const response = await axios.post("http://localhost:8089/delete_task", {
      taskId,
    });
    alert(`任务 ${taskId} 删除成功！`);
    // 可选：刷新任务列表
    getAllTasks();
  } catch (error) {
    console.error("删除任务失败:", error);
    alert("删除任务失败，请稍后重试！");
  }
}
const finishTask = async (taskId) => {
  const task = tasks.value.find(t => t.ID === taskId);
  try {
    const response = await axios.post("http://localhost:8089/finish_task", {
      taskId,
    });
    alert(`任务 ${taskId} 删除成功！`);
    // 可选：刷新任务列表
    getAllTasks();
  } catch (error) {
    console.error("删除任务失败:", error);
    alert("删除任务失败，请稍后重试！");
  }
}
// 新增变量：控制下一轮弹窗显示和输入的根模型 ID
const showNextRoundModal = ref(false);
const nextRoundRootModelId = ref('');
const selectedTaskId = ref(null);

// 修改进入下一轮任务逻辑
const startNextRound = (taskId) => {
  selectedTaskId.value = taskId; // 确保正确设置任务 ID
  showNextRoundModal.value = true; // 打开弹窗
};

const submitNextRound = async () => {
  if (!nextRoundRootModelId.value) {
    alert("根模型 ID 不能为空！");
    return;
  }

  try {
    console.log({
      taskId: selectedTaskId.value, // 确保传递正确的任务 ID
      rootModelId: nextRoundRootModelId.value,
      username: userInfo.value.username,
    })
    const response = await axios.post("http://localhost:8089/next_task_round", {
      taskId: selectedTaskId.value, // 确保传递正确的任务 ID
      rootModelId: nextRoundRootModelId.value,
      username: userInfo.value.username,
    });
    alert(`任务 ${selectedTaskId.value} 已进入下一轮！`);
    // 清空输入框和任务 ID，并关闭弹窗
    nextRoundRootModelId.value = '';
    selectedTaskId.value = null;
    showNextRoundModal.value = false;
    // 刷新任务列表
    getAllTasks();
  } catch (error) {
    console.error("进入下一轮失败:", error);
    alert("进入下一轮失败，请稍后重试！");
  }
};

const cancelNextRound = () => {
  // 清空任务 ID 并关闭弹窗
  selectedTaskId.value = null;
  showNextRoundModal.value = false;
};

// 新任务弹窗相关变量
const showNewTaskModal = ref(false);
const newTaskBonus = ref('');
const newTaskRootModelId = ref('');

// 提交新任务
const submitNewTask = async () => {
  if (newTaskBonus.value < 0) {
    alert("代币奖励值不能为负数！");
    return;
  }

  try {
    const response = await axios.post("http://localhost:8089/new_task", {
      username: userInfo.value.username,
      bonus: newTaskBonus.value,
      rootModelId: newTaskRootModelId.value,
    });
    alert("新任务发布成功！");
    // 清空输入框并关闭弹窗
    newTaskBonus.value = '';
    newTaskRootModelId.value = '';
    showNewTaskModal.value = false;
    // 刷新任务列表
    getAllTasks();
  } catch (error) {
    console.error("发布新任务失败:", error);
    alert("发布新任务失败，请稍后重试！");
  }
};

// 页面加载时获取任务列表
onMounted(() => {
  getAllTasks()
})

// 退出登录
const handleLogout = async () => {
  try {
    const response = await axios.post("http://localhost:8089/log_out", {
      username: userInfo.value.username,
      organization: userInfo.value.organization, // 传递用户组织信息
    });
    console.log("用户组织:", response.data.organization); // 处理返回的用户组织信息
  } catch (error) {
    console.error("注销请求失败:", error);
  }
  // 清除用户状态
  localStorage.removeItem('authToken');
  localStorage.removeItem('userInfo');

  // 跳转到登录页并阻止返回
  router.replace('/').then(() => {
    window.location.reload(); // 可选：完全重置应用状态
  });
};
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

/* 新增模块样式 */
.my-tasks-container {
  margin-top: 20px;
  background: #f9f9f9;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.my-tasks-container h3 {
  color: #2c3e50;
  margin-bottom: 10px;
}

/* 新增按钮样式 */
.new-task-button {
  margin-top: 20px;
  padding: 10px 20px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.new-task-button:hover {
  background-color: #2980b9;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 10px;
  width: 400px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
}

.modal-content h3 {
  margin-bottom: 20px;
  color: #2c3e50;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: #34495e;
}

.form-group input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 5px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
}

.submit-button {
  padding: 10px 20px;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.submit-button:hover {
  background-color: #45a049;
}

.cancel-button {
  padding: 10px 20px;
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.cancel-button:hover {
  background-color: #c0392b;
}
</style>