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
            <div v-if="userInfo" class="welcome-container">
                <h1>个人信息界面</h1>
                <ul class="system-info">
                    <li>👤 用户名：{{ username_now }}</li>
                    <li>📊 所属组织：{{ organization_now }}</li>
                    <li>🔑 公钥哈希：{{ pubkey_now }}</li>
                    <li>🪙 持有积分：{{ token }}</li>
                </ul>
                <button class="upload-key-btn" @click="showUploadPanel = true">上传公钥</button>
            </div>

            <!-- 已上传模型展示 -->
            <div v-if="postedModels.length > 0" class="models-container">
                <h2>已上传模型</h2>
                <table class="models-table">
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>模型ID</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(model, index) in postedModels" :key="index">
                            <td>{{ index + 1 }}</td>
                            <td>{{ model }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div v-else class="no-models-container">
                <p>暂无已上传模型。</p>
            </div>

            <!-- 已创建任务展示 -->
            <div v-if="tasks.length > 0" class="tasks-container">
                <h2>已创建任务</h2>
                <ul class="tasks-list">
                    <li v-for="(task, index) in tasks" :key="index">
                        📝 任务名称：{{ task }}
                    </li>
                </ul>
            </div>
            <div v-else class="no-tasks-container">
                <p>暂无已创建任务。</p>
            </div>



            <!-- 上传公钥面板 -->
            <div v-if="showUploadPanel" class="upload-panel">
                <h2>上传公钥</h2>
                <textarea v-model="publicKeyInput" placeholder="请输入公钥"></textarea>
                <div class="panel-actions">
                    <button @click="submitPublicKey">提交</button>
                    <button @click="showUploadPanel = false">取消</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

// 定义响应式变量
const username_now = ref("加载中...")
const organization_now = ref("加载中...")
const pubkey_now = ref("加载中...")
const token = ref("加载中...")
const showUploadPanel = ref(false) // 控制上传公钥面板的显示
const publicKeyInput = ref("") // 存储用户输入的公钥
const postedModels = ref([]); // 存储已上传模型的数组，直接展示内容
const tasks = ref([]); // 存储已创建任务的数组

const pubkeycommit_form = ref({
    username: "",
    publickey: "",
});

// 从 localStorage 获取用户信息
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

// 自动获取用户信息的方法
const getinfo = async () => {
    try {
        // 向后端发送请求
        const response = await axios.post("http://localhost:8089/get_user_info", {
            username: userInfo.value.username,
        });

        // 更新响应式变量
        username_now.value = response.data.user.username;
        organization_now.value = response.data.user.organization;
        pubkey_now.value = response.data.user.pubkeyhash;
        token.value = response.data.user.token;
        postedModels.value = response.data.user.posted || []; // 更新已上传模型数组
        tasks.value = response.data.user.tasks || []; // 更新已创建任务数组
        console.log("已创建任务:", tasks.value);
    } catch (error) {
        console.error("获取用户信息失败:", error);
        alert("获取用户信息失败，请稍后重试！");
    }
};

// 提交公钥的方法
const submitPublicKey = async () => {
    try {
        // 向后端发送公钥
        console.log("上传公钥的响应:", publicKeyInput.value);
        const response = await axios.post("http://localhost:8089/upload_public_key", {
            username: userInfo.value.username,
            Pubkeyhash: publicKeyInput.value,
        });
        // 更新公钥哈希
        pubkey_now.value = response.data.pubkeyhash;

        alert("公钥上传成功！");
        showUploadPanel.value = false; // 关闭面板
    } catch (error) {
        console.error("上传公钥失败:", error);
        alert("上传公钥失败，请稍后重试！");
    }
};

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

// 组件加载时自动调用 getinfo
onMounted(() => {
    getinfo();
});

// 模拟查看模型详情的方法
const viewModelDetails = (model) => {
    alert(`查看模型详情: ${model.name}`);
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

/* 加载中提示样式 */
.loading-container {
    text-align: center;
    padding: 20px;
    font-size: 1.2rem;
    color: #7f8c8d;
}

/* 上传公钥面板样式 */
.upload-panel {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
    width: 400px;
    z-index: 1000;
}

.upload-panel h2 {
    margin-bottom: 20px;
    text-align: center;
}

.upload-panel textarea {
    width: 100%;
    height: 100px;
    margin-bottom: 20px;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 5px;
    resize: none;
}

.panel-actions {
    display: flex;
    justify-content: space-between;
}

.panel-actions button {
    padding: 10px 20px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    transition: background-color 0.3s ease;
}

.panel-actions button:first-child {
    background-color: #4caf50;
    color: white;
}

.panel-actions button:first-child:hover {
    background-color: #45a049;
}

.panel-actions button:last-child {
    background-color: #e74c3c;
    color: white;
}

.panel-actions button:last-child:hover {
    background-color: #c0392b;
}

/* 已上传模型展示样式 */
.models-container {
    margin-top: 20px;
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    max-height: 200px; /* 固定高度 */
    overflow-y: auto; /* 添加滚动条 */
}

.models-container h2 {
    margin-bottom: 15px;
    color: #2c3e50;
}

/* 模型表格样式 */
.models-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 10px;
}

.models-table th, .models-table td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: center;
}

.models-table th {
    background-color: #f4f4f4;
    color: #333;
    font-weight: bold;
}

.models-table tr:nth-child(even) {
    background-color: #f9f9f9;
}

.models-table tr:hover {
    background-color: #f1f1f1;
}
</style>