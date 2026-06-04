<!--
/**
 * menu.vue - 顶部导航菜单组件
 *
 * 功能说明:
 * 这是应用顶部的导航菜单组件，显示登录/用户信息。
 * 根据用户登录状态显示不同的内容：
 * - 未登录：显示"登录"按钮
 * - 已登录：显示用户头像和退出选项
 *
 * 技术栈:
 * - Vue 3 Composition API (setup 语法糖)
 * - Element Plus UI 组件
 * - JWT Token 解析
 */
-->

<!-- ============================================================
     <template> 区域
     ============================================================ -->
<template>
  <!--
    el-menu - Element Plus 的导航菜单组件
    mode="horizontal" - 水平模式（横向显示）
    background-color - 背景色（深灰色）
    text-color - 文字颜色（白色）
    menu-trigger="click" - 点击触发子菜单
    @select - 菜单选中事件
  -->
  <el-menu
    class="el-menu-demo"
    mode="horizontal"
    :ellipsis="false"
    background-color="#545c64"
    text-color="#ffffff"
    menu-trigger="click"
    @select="handleSelect">

    <!-- .flex-grow - 占位元素，将后续菜单推到右侧 -->
    <div class="flex-grow" />

    <!-- --------------------------------------------------------
         登录按钮 - 仅在用户未登录时显示
         v-if / :style 条件渲染
         userInfo.user_id == 0 表示未登录
         -------------------------------------------------------- -->
    <el-menu-item index="0" :style="{display: userInfo.user_id != 0 ?'none':''}">
      <!-- 登录按钮 -->
      <el-button>登录</el-button>
    </el-menu-item>

    <!-- --------------------------------------------------------
         用户信息 - 仅在用户已登录时显示
         userInfo.user_id != 0 表示已登录
         -------------------------------------------------------- -->
    <el-sub-menu index="1" :style="{display: userInfo.user_id == 0 ?'none':''}">
      <!-- #title 是具名插槽，显示在触发器上 -->
      <template #title>
        <!-- el-avatar - 用户头像组件 -->
        <!-- :src 绑定头像 URL，如果加载失败会显示默认占位符 -->
        <el-avatar :src="userInfo.avatar"></el-avatar>
      </template>

      <!-- 退出菜单项 -->
      <el-menu-item index="1-1">退出</el-menu-item>
    </el-sub-menu>

  </el-menu>
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

// ref - 创建响应式数据的函数
// onBeforeMount - 组件挂载前的生命周期钩子
import { ref, onBeforeMount } from 'vue'

// getCookie - 从 Cookie 获取数据的工具函数
// deleteCookie - 删除 Cookie 的工具函数
import { getCookie, deleteCookie } from '../utils/utils.ts'

// ============================================================
// 2. 定义响应式状态
// ============================================================

/**
 * userInfo - 用户信息响应式对象
 *
 * 包含用户的基本信息：
 * - name: 用户名
 * - user_id: 用户 ID（0 表示未登录，非 0 表示已登录）
 * - avatar: 用户头像 URL
 */
let userInfo = ref<{ name: string; avatar: string; user_id: number }>({
  name: "",       // 用户名，初始为空
  user_id: 0,     // 用户 ID，0 表示未登录状态
  avatar: "",     // 头像 URL，初始为空
});

// ============================================================
// 3. 组件挂载前生命周期钩子
// ============================================================

/**
 * onBeforeMount - 组件挂载到 DOM 前执行
 *
 * 用于初始化用户信息：从 Cookie 中读取 JWT Token 并解析
 */
onBeforeMount(() => {
  // --------------------------------------------------------
  // 从 Cookie 获取访问令牌
  // --------------------------------------------------------
  // sso_0voice_access_token 是登录后后端设置的 JWT Token
  // 存储在 Cookie 中，可用于身份验证
  let access_token = getCookie("sso_0voice_access_token");

  console.log(access_token);  // 调试：打印 Token

  // --------------------------------------------------------
  // 解析 JWT Token
  // --------------------------------------------------------
  // JWT 格式：header.payload.signature（用 . 分隔）
  // 例如：eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4ifQ.xxx
  //
  // 我们只需要解析 payload（第二部分），包含用户信息

  // split(".") - 用点号分割成数组
  // [0] = header, [1] = payload, [2] = signature
  let list = access_token?.split(".");

  // --------------------------------------------------------
  // 解析用户信息
  // --------------------------------------------------------
  if (list) {
    // --------------------------------------------------------
    // Base64 解码 payload
    // --------------------------------------------------------
    // atob() - 将 Base64 编码的字符串解码为原始字符串
    // btoa() - 将字符串编码为 Base64（编码时使用）
    //
    // JWT payload 示例：
    // {"user_id":123,"name":"张三","avatar":"http://...","exp":1625123456}
    // 这个 JSON 字符串经过 Base64 编码后存储在 Token 的第二部分

    // JSON.parse() - 将 JSON 字符串解析为 JavaScript 对象
    userInfo.value = JSON.parse(atob(list[1]));
  }
});

// ============================================================
// 4. 菜单选择处理函数
// ============================================================

/**
 * handleSelect - 处理菜单项点击
 *
 * @param key - 被选中的菜单项 index
 * @param keyPath - 菜单项的路径数组
 *
 * 通过 switch 语句处理不同的菜单操作
 */
const handleSelect = (key: string, keyPath: string[]) => {
  console.log(key, keyPath);  // 调试：打印选中的菜单信息

  // 根据选中的菜单项执行不同操作
  switch (key) {
    // --------------------------------------------------------
    // case '0' - 登录按钮
    // --------------------------------------------------------
    case '0':
      // 跳转到用户中心/登录页面
      // VITE_USER_CENTER 是 Vite 环境变量，存储用户中心的 URL
      window.location.href = import.meta.env.VITE_USER_CENTER;
      break;

    // --------------------------------------------------------
    // case '1-1' - 退出菜单项
    // --------------------------------------------------------
    case '1-1':
      // 删除认证 Cookie
      deleteCookie("sso_0voice_access_token");

      // 刷新当前页面
      // 刷新后页面会重新检查登录状态（userInfo 会被重新初始化）
      window.location.href = window.location.href;
      break;

    // --------------------------------------------------------
    // default - 其他未处理的情况
    // --------------------------------------------------------
    default:
      // 不做任何处理
      break;
  }
};
</script>

<!-- ============================================================
     <style> 区域
     ============================================================ -->
<style>
/* --------------------------------------------------------
   .flex-grow 类样式
   flex-grow: 1 使元素占据所有可用空间
   这会将后续元素推到容器的另一侧
   -------------------------------------------------------- */
.flex-grow {
  flex-grow: 1;
}
</style>
