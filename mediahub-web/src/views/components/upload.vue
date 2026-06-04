<!--
/**
 * upload.vue - 文件上传组件
 *
 * 功能说明:
 * 提供文件上传功能，包含：
 * 1. 上传按钮 - 点击后打开文件选择器
 * 2. URL 文本框 - 显示上传后的文件地址
 * 3. 复制按钮 - 一键复制文件 URL
 *
 * 技术栈:
 * - Vue 3 Composition API
 * - Element Plus 组件库
 * - Clipboard API
 */
-->

<!-- ============================================================
     <template> 区域
     ============================================================ -->
<template>
  <div>
    <!-- --------------------------------------------------------
         上传容器
         -------------------------------------------------------- -->
    <!-- 外层容器：白色背景，圆角边框 -->
    <div style="background-color: white;width:55.375rem;height:2.875rem;padding:0.4375rem;border-radius: 5px">
      <!-- 内层容器：相对定位作为子元素的定位参考 -->
      <div style="position: relative;width:100%;height:100%">

        <!-- --------------------------------------------------------
             上传按钮
             -------------------------------------------------------- -->
        <!-- el-button - Element Plus 按钮组件 -->
        <!-- type="success" 绿色样式，表示主要操作 -->
        <!-- @click="handleUpload" 点击事件处理 -->
        <el-button
          type="success"
          class="grid-content ep-bg-purple"
          style="position:absolute;left:0;height:2.875rem;width:4.875rem;font-family:Microsoft YaHei;font-size: 1.25rem;"
          @click="handleUpload">
          上传
        </el-button>

        <!-- --------------------------------------------------------
             URL 文本框
             -------------------------------------------------------- -->
        <!-- 隐藏的文件输入框（用于选择文件） -->
        <!-- v-model 绑定 URL 数据，显示上传结果 -->
        <input
          style="width:40.5rem;height:2.75rem;position:absolute;left:5.8rem;border:none;padding:0.0625rem;font-size: 1.5rem;color:#606266;outline:0;"
          id="upload-input"
          v-model="data.url.value" />

        <!-- --------------------------------------------------------
             复制按钮
             -------------------------------------------------------- -->
        <!-- @click="handleCopy" 点击复制 URL -->
        <el-button
          type="success"
          class="grid-content ep-bg-purple"
          style="position:absolute;right:0;height:2.875rem;width:8.0625rem;font-family:Microsoft YaHei;font-size: 1.25rem;"
          @click="handleCopy">
          复制地址
        </el-button>

      </div>
    </div>
  </div>
</template>

<!-- ============================================================
     <script setup> 区域
     ============================================================ -->
<script lang="ts" setup>
/**
 * lang="ts" - 使用 TypeScript
 */

// ============================================================
// 1. 导入依赖
// ============================================================

// ref - 创建响应式数据
import { ref } from 'vue'

// uploadFile - 文件上传 API 函数
import { uploadFile } from '../../api/api'

// ElMessage - 消息提示组件
// ElLoading - 加载指示器组件
import { ElMessage, ElLoading } from 'element-plus'

// ============================================================
// 2. 定义响应式数据
// ============================================================

/**
 * data - 组件数据对象
 *
 * url - 上传成功后的文件 URL
 */
const data = {
  // 初始为空字符串，上传成功后更新
  url: ref("")
};

// ============================================================
// 3. 定义响应数据类型
// ============================================================

/**
 * fileUploadRes - 文件上传响应数据类型
 *
 * 使用 ES6 Class 定义，匹配后端返回的数据结构
 */
class fileUploadRes {
  // 上传成功后的文件 URL
  url: string;

  // 构造函数
  constructor(url: string) {
    this.url = url;
  }
}

// ============================================================
// 4. 上传处理函数
// ============================================================

/**
 * handleUpload - 处理文件上传
 *
 * 点击上传按钮后：
 * 1. 创建隐藏的文件输入框
 * 2. 打开文件选择器
 * 3. 选择文件后调用上传 API
 * 4. 显示加载状态
 * 5. 上传完成后显示结果
 */
function handleUpload() {
  console.log("handleUpload");  // 调试日志

  // --------------------------------------------------------
  // 步骤 1: 创建隐藏的文件输入框
  // --------------------------------------------------------
  // document.createElement('input') 创建新的 DOM 元素
  const input = document.createElement('input');

  // 设置 input 类型为 file
  input.type = "file";

  // 设置唯一 ID
  input.id = "file-upload";

  // --------------------------------------------------------
  // 步骤 2: 添加文件选择事件监听器
  // --------------------------------------------------------
  // addEventListener 监听 'change' 事件
  // 当用户选择文件后触发
  input.addEventListener("change", (event) => {
    // event.target 是触发事件的元素（input）
    // .files 是 FileList 对象，包含用户选择的文件列表
    const files = (event.target as HTMLInputElement).files;

    // --------------------------------------------------------
    // 步骤 3: 检查是否有文件被选中
    // --------------------------------------------------------
    // files 存在且长度大于 0，表示用户选择了文件
    if (files && files.length > 0) {
      // --------------------------------------------------------
      // 步骤 4: 显示加载指示器
      // --------------------------------------------------------
      // ElLoading.service() 创建全屏加载遮罩
      // lock: true 锁定页面，防止用户操作
      // text 显示加载提示文字
      // background 设置遮罩背景色（半透明黑色）
      const loading = ElLoading.service({
        lock: true,
        text: "文件上传中。。。",
        background: "rgba(0,0,0,0.7)"
      });

      // --------------------------------------------------------
      // 步骤 5: 构造 FormData
      // --------------------------------------------------------
      // FormData 是浏览器提供的 API
      // 用于构造 multipart/form-data 格式的请求体
      const formData = new FormData();

      // append() 添加字段
      // 'file' 是后端接收文件的字段名
      // files[0] 是用户选择的第一个文件
      formData.append("file", files[0]);

      // --------------------------------------------------------
      // 步骤 6: 调用上传 API
      // --------------------------------------------------------
      uploadFile<fileUploadRes>({ formData: formData }).then(function (res) {
        // 上传成功
        // res.data?.url 使用可选链，安全访问返回的 URL
        // || "" 处理空值情况
        data.url.value = res.data?.url || "";

      }).catch(function (res) {
        // 上传失败
        console.log(res.message);  // 打印错误信息

      }).finally(function () {
        // 无论成功或失败，最后都要关闭加载指示器
        loading.close();
      });
    }
  });

  // --------------------------------------------------------
  // 步骤 7: 触发文件选择对话框
  // --------------------------------------------------------
  // .click() 方法模拟点击，打开文件选择器
  input.click();
}

// ============================================================
// 5. 复制处理函数
// ============================================================

/**
 * handleCopy - 处理复制 URL 到剪贴板
 *
 * 使用浏览器 Clipboard API 复制文本
 */
function handleCopy() {
  // --------------------------------------------------------
  // 使用 Clipboard API 复制文本
  // --------------------------------------------------------
  // navigator.clipboard.writeText() 将文本写入系统剪贴板
  // 返回 Promise，支持 async/await
  navigator.clipboard.writeText(data.url.value).then(function () {
    // 复制成功
    // ElMessage 显示成功提示
    ElMessage({
      message: "已复制到剪切板",
      type: "success",
    });

  }).catch(function () {
    // 复制失败（可能是浏览器不支持或权限问题）
    // ElMessage.error 显示错误提示
    ElMessage.error("复制失败，请手动复制文本框内的链接");
  });
}
</script>
