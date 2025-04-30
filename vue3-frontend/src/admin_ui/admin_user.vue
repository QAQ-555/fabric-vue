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
      <div v-if="$route.path === '/admin/main'" class="welcome-container">
        <h1>欢迎回来，{{ userInfo.username }}！</h1>
        <p>组织：{{ userInfo.organization }}</p>
        <p>当前系统信息：</p>
        <ul class="system-info">
          <li>🕒 登录时间：{{ loginTime }}</li>
          <li>📌 当前版本：v7.38d</li>
        </ul>
      </div>

      <!-- 用户管理 -->
      <div class="content-placeholder" v-else>
        <h2>用户管理模块</h2>

        <!-- 已验证用户 -->
        <div v-if="verifiedUsers.length > 0" class="users-container">
          <h3>已验证用户</h3>
          <table class="users-table">
            <thead>
              <tr>
                <th>#</th>
                <th>用户名</th>
                <th>组织</th>
                <th>公钥哈希</th>
                <th>令牌</th>
                <th>是否管理员</th>
                <th>是否接受</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(user, index) in verifiedUsers" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ user.username }}</td>
                <td>{{ user.organization }}</td>
                <td>{{ user.pubkeyhash }}</td>
                <td>{{ user.token }}</td>
                <td>{{ user.isAdmin ? "是" : "否" }}</td>
                <td>{{ user.isAccepted ? "是" : "否" }}</td>
                <td>
                  <button class="delete-button" @click="openDeleteDialog(user)">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else>
          <p>暂无已验证用户。</p>
        </div>

        <!-- 未验证用户 -->
        <div v-if="unverifiedUsers.length > 0" class="users-container">
          <h3>未验证用户</h3>
          <table class="users-table">
            <thead>
              <tr>
                <th>#</th>
                <th>用户名</th>
                <th>组织</th>
                <th>公钥哈希</th>
                <th>令牌</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(user, index) in unverifiedUsers" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ user.username }}</td>
                <td>{{ user.organization }}</td>
                <td>{{ user.pubkeyhash }}</td>
                <td>{{ user.token }}</td>
                <td>
                  <button @click="openManageDialog(user)">管理</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else>
          <p>暂无未验证用户。</p>
        </div>
      </div>

      <!-- 管理弹窗 -->
      <div v-if="showDialog" class="dialog-overlay">
        <div class="dialog">
          <h3>管理用户: {{ selectedUser.username }}</h3>
          <div class="dialog-content">
            <div class="dialog-section">
              <h4>接受状态</h4>
              <div class="radio-group">
                <label class="radio-label">
                  <input
                    type="radio"
                    v-model="isAccepted_var"
                    :value="true"
                  />
                  接受
                </label>
                <label class="radio-label">
                  <input
                    type="radio"
                    v-model="isAccepted_var"
                    :value="false"
                  />
                  不接受
                </label>
              </div>
            </div>
            <div
              class="dialog-section"
              :class="{ disabled: !isAccepted_var }"
            >
              <h4>用户角色</h4>
              <div class="radio-group">
                <label class="radio-label">
                  <input
                    type="radio"
                    v-model="isAdmin_var"
                    :disabled="!isAccepted_var"
                    :value="true"
                  />
                  管理员
                </label>
                <label class="radio-label">
                  <input
                    type="radio"
                    v-model="isAdmin_var"
                    :disabled="!isAccepted_var"
                    :value="false"
                  />
                  普通用户
                </label>
              </div>
            </div>
          </div>
          <div class="dialog-actions">
            <button class="save-button" @click="verifyUser">保存</button>
            <button class="cancel-button" @click="closeDialog">取消</button>
          </div>
        </div>
      </div>

      <!-- 删除确认弹窗 -->
      <div v-if="showDeleteDialog" class="dialog-overlay">
        <div class="dialog">
          <h3>确认删除用户</h3>
          <p>
            此操作不可逆！
            请输入目标用户名以确认删除：
          </p>
          <input
            v-model="deleteConfirmation"
            type="text"
            placeholder="输入用户名确认"
            class="input-confirmation"
          />
          <div class="dialog-actions">
            <button
              class="save-button"
              :disabled="deleteConfirmation !== selectedUser.username"
              @click="confirmDeleteUser"
            >
              确认删除
            </button>
            <button class="cancel-button" @click="closeDeleteDialog">
              取消
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 引入依赖
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";

const router = useRouter();
const userInfo = ref(JSON.parse(localStorage.getItem("userInfo") || "{}"));
const loginTime = ref(new Date().toLocaleString());
const users = ref([]);
const isAccepted_var =ref(false);
const isAdmin_var =ref(false);
// 已验证和未验证用户的分类
const verifiedUsers = ref([]);
const unverifiedUsers = ref([]);

// 获取用户列表并分类
const getAllUsers = async () => {
  try {
    const response = await axios.post("http://localhost:8089/get_all_users");
    const allUsers = response.data.users || [];
    verifiedUsers.value = allUsers.filter((user) => user.isVerified);
    unverifiedUsers.value = allUsers.filter((user) => !user.isVerified);
    console.log("获取的用户列表:", response.data);
  } catch (error) {
    console.error("获取用户失败:", error);
    alert("获取用户失败，请稍后重试！");
  }
};

// 删除用户
const deleteUser = async (username) => {
  const confirmation = prompt(
    `此操作不可逆！请输入用户名 "${username}" 以确认删除：`
  );
  if (confirmation !== username) {
    alert("用户名输入错误，操作已取消！");
    return;
  }

  try {
    await axios.post("http://localhost:8089/delete_user", { username });
    alert(`用户 ${username} 已成功删除！`);
    getAllUsers(); // 刷新用户列表
  } catch (error) {
    console.error(`删除用户 ${username} 失败:`, error);
    alert("删除用户失败，请稍后重试！");
  }
};

// 验证用户
const verifyUser = async () => {
  try {
    console.log({
      username: selectedUser.value.username,
      isAdmin: isAdmin_var.value,
      isAccepted: isAccepted_var.value,
    });
    await axios.post("http://localhost:8089/verify_user", {
      username: selectedUser.value.username,
      isAdmin: isAdmin_var.value,
      isAccepted: isAccepted_var.value,
    });
    alert(`用户 ${selectedUser.value.username} 状态已成功更新！`);
    getAllUsers(); // 刷新用户列表
    closeDialog(); // 关闭弹窗
  } catch (error) {
    console.error(`验证用户 ${selectedUser.value.username} 失败:`, error);
    alert("验证用户失败，请稍后重试！");
  }
};

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
    window.location.reload();
  });
};

// 页面加载时获取用户列表
onMounted(() => {
  getAllUsers();
});

const showDialog = ref(false);
const selectedUser = ref({});

// 打开管理弹窗
const openManageDialog = (user) => {
  selectedUser.value = { ...user };
  isAccepted_var.value = user.isAccepted || false;
  isAdmin_var.value = user.isAdmin || false;
  showDialog.value = true;
};

// 关闭管理弹窗
const closeDialog = () => {
  showDialog.value = false;
  selectedUser.value = {};
};

const showDeleteDialog = ref(false);
const deleteConfirmation = ref("");

// 打开删除确认弹窗
const openDeleteDialog = (user) => {
  selectedUser.value = { ...user };
  deleteConfirmation.value = "";
  showDeleteDialog.value = true;
};

// 关闭删除确认弹窗
const closeDeleteDialog = () => {
  showDeleteDialog.value = false;
  selectedUser.value = {};
};

// 确认删除用户
const confirmDeleteUser = async () => {
  await deleteUser(selectedUser.value.username);
  closeDeleteDialog();
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
  width: calc(100% - 40px);
  /* 适配侧边栏padding */
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

/* 用户表格样式 */
.users-container {
  margin-top: 20px;
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  max-height: 400px;
  /* 固定高度 */
  overflow-y: auto;
  /* 添加滚动条 */
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

.users-table th,
.users-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: center;
}

.users-table th {
  background-color: #f4f4f4;
  color: #333;
  font-weight: bold;
}

.users-table tr:nth-child(even) {
  background-color: #f9f9f9;
}

.users-table tr:hover {
  background-color: #f1f1f1;
}

.users-table button {
  padding: 5px 10px;
  border: none;
  background-color: #3498db;
  color: white;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.users-table button:hover {
  background-color: #2980b9;
}

/* 删除按钮样式 */
.delete-button {
  padding: 5px 10px;
  border: none;
  background-color: #e74c3c;
  color: white;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.delete-button:hover {
  background-color: #c0392b;
  transform: scale(1.05);
}

.delete-button:active {
  transform: scale(0.95);
}

/* 弹窗样式 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.dialog {
  background: white;
  padding: 20px;
  border-radius: 10px;
  width: 300px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.dialog-content {
  margin-bottom: 20px;
}

.dialog-section {
  margin-bottom: 20px;
}

.dialog-section.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.dialog-section h4 {
  margin-bottom: 10px;
  font-size: 1rem;
  color: #2c3e50;
}

.radio-group {
  display: flex;
  gap: 20px;
}

.radio-label {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  color: #34495e;
}

.radio-label input {
  margin-right: 5px;
}

.dialog-actions {
  display: flex;
  justify-content: space-between;
}

.save-button {
  padding: 10px 20px;
  background-color: #27ae60;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.save-button:hover {
  background-color: #219150;
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

.input-confirmation {
  width: 90%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 1rem;
}
</style>
