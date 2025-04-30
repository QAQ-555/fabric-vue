<template>
  <div class="system-container">
    <!-- 左侧导航栏 -->
    <div class="sidebar">
      <div class="logo">AI System</div>
      <nav>
        <!-- 导航链接 -->
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

      <!-- 模型管理 -->
      <div class="content-placeholder" v-else>
        <h2>模型管理</h2>
        <!-- 上传模型 -->
        <div class="button-group">
          <textarea
            v-model="modelSignature"
            placeholder="请输入模型签名"
            class="signature-textarea"
          ></textarea>
          <div class="button-container">
            <button class="action-button" @click="handleUpload">选择并上传模型</button>
            <input type="file" ref="fileInput" style="display: none" @change="uploadFile" />
          </div>
        </div>

        <!-- 下载模型 -->
        <div class="download-section">
          <h3>下载模型</h3>
          <p>请输入模型的 CID 以下载文件。</p>
          <button class="download-button" @click="handleDownload">下载模型</button>
        </div>

        <!-- 查询模型 CID -->
        <div class="query-section">
          <h3>查询模型 CID</h3>
          <p>请输入模型 ID：</p>
          <input v-model="modelId" placeholder="模型 ID" class="query-input" />
          <button class="query-button" @click="queryModelCid">查询 CID</button>
          <p v-if="modelCid">模型 CID：{{ modelCid }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 引入依赖
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { create } from 'ipfs-http-client'

// 初始化变量
const router = useRouter()
const ipfs = create({ url: 'http://localhost:5001' })
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
const loginTime = ref(new Date().toLocaleString())
const pendingTasks = ref(3)
const modelSignature = ref('')
const fileInput = ref(null)

// 新增变量
const modelId = ref('')
const modelCid = ref('')

// 上传模型方法
const handleUpload = () => {
  if (!modelSignature.value.trim()) {
    alert('请先输入模型签名')
    return
  }
  fileInput.value.click()
}

const uploadFile = async (event) => {
  const file = event.target.files[0]
  if (!file) {
    console.log('未选择文件')
    return
  }

  try {
    console.log('开始上传文件到 IPFS...')
    const added = await ipfs.add(file)
    console.log(`文件上传成功，CID: ${added.cid}`)

    await sendToBackend(modelSignature.value, added.cid.toString())
    alert(`文件已上传到 IPFS，CID: ${added.cid}`)
  } catch (error) {
    console.error('上传失败:', error)
    alert('上传失败，请重试')
  }
}

const sendToBackend = async (signature, cid) => {
  try {
    const response = await axios.post('http://localhost:8089/upload_model', {
      username: userInfo.value.username,
      signature,
      cid,
    })
    console.log('后端返回:', response.data)
    alert('模型签名和 CID 已成功发送到后端！')
  } catch (error) {
    console.error('发送到后端失败:', error)
    alert('发送到后端失败，请稍后重试！')
  }
}

// 下载模型方法
const handleDownload = async () => {
  const cid = prompt('请输入模型的 CID：')
  if (!cid) {
    alert('CID 不能为空')
    return
  }

  try {
    console.log(`开始从 IPFS 下载文件，CID: ${cid}`)
    const stream = ipfs.cat(cid)
    const chunks = []

    for await (const chunk of stream) {
      chunks.push(chunk)
    }

    const fileContent = new Blob(chunks)
    const url = window.URL.createObjectURL(fileContent)
    const a = document.createElement('a')
    a.href = url
    a.download = `${cid}`
    a.click()
    window.URL.revokeObjectURL(url)

    alert('文件已成功下载！')
  } catch (error) {
    console.error('下载失败:', error)
    alert('下载失败，请重试')
  }
}

// 查询模型 CID 方法
const queryModelCid = async () => {
  if (!modelId.value.trim()) {
    alert('模型 ID 不能为空')
    return
  }

  try {
    const response = await axios.post('http://localhost:8089/query_model_cid', {
      modelId: modelId.value,
    })
    modelCid.value = response.data.cid
    alert('模型 CID 查询成功！')
  } catch (error) {
    console.error('查询模型 CID 失败:', error)
    alert('查询模型 CID 失败，请稍后重试！')
  }
}

// 退出登录方法
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

/* 按钮组样式 */
.button-group {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.button-container {
  display: flex;
  justify-content: center;
  width: 100%;
}

.action-button {
  padding: 10px 20px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.action-button:hover {
  background-color: #2980b9;
}

.signature-textarea {
  width: 100%;
  height: 100px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 14px;
  resize: none;
  margin-bottom: 10px;
}

/* 下载按钮区域样式 */
.download-section {
  margin-top: 30px;
  padding: 20px;
  background-color: #ecf0f1;
  border-radius: 10px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.download-section h3 {
  color: #2c3e50;
  margin-bottom: 10px;
}

.download-section p {
  color: #7f8c8d;
  margin-bottom: 20px;
}

.download-button {
  padding: 10px 20px;
  background-color: #e67e22;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.download-button:hover {
  background-color: #d35400;
}

/* 查询模型 CID 样式 */
.query-section {
  margin-top: 30px;
  padding: 20px;
  background-color: #ecf0f1;
  border-radius: 10px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.query-section h3 {
  color: #2c3e50;
  margin-bottom: 10px;
}

.query-section p {
  color: #7f8c8d;
  margin-bottom: 10px;
}

.query-input {
  width: 80%;
  padding: 10px;
  margin-bottom: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.query-button {
  padding: 10px 20px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.query-button:hover {
  background-color: #2980b9;
}
</style>