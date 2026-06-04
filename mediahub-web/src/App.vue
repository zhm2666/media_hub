<!--
/**
 * App.vue - Vue 应用根组件
 *
 * 功能说明:
 * 这是 Vue 应用的根组件，作为所有页面的容器。
 * 目前主要用于渲染首页组件。
 *
 * 组件结构:
 * <script setup> - 组件逻辑区域，使用 Composition API
 * <template> - 组件模板区域，定义 HTML 结构
 * <style scoped> - 组件样式区域，scoped 表示样式仅在当前组件内生效
 */
-->

<!-- ============================================================
     <script setup> 区域
     使用 TypeScript 和 Composition API 定义组件逻辑
     ============================================================ -->
<script setup lang="ts">
/**
 * lang="ts" - 声明使用 TypeScript 语言
 * setup 是 Vue 3 Composition API 的入口点
 * 所有在 setup 中定义的变量和函数都可以在模板中直接使用
 */

// ============================================================
// 1. 导入首页组件
// ============================================================
// 从 views/home.vue 导入首页组件
// 首页包含轮播图、图片展示区和上传功能
import home from './views/home.vue'
</script>

<!-- ============================================================
     <template> 区域
     定义组件的 HTML 模板结构
     ============================================================ -->
<template>
  <!-- home 组件 - 首页主组件
    所有页面内容都被包含在首页组件中 -->
  <home></home>
</template>

<!-- ============================================================
     <style scoped> 区域
     定义组件的私有样式
     scoped 属性确保这些样式只影响当前组件
     ============================================================ -->
<style scoped>
/**
 * .logo 类 - 应用于 Logo 的样式
 */

/* logo 容器高度设为 6em（相对于当前字体大小的 6 倍） */
.logo {
  height: 6em;
  padding: 1.5em;  /* 内边距，给 logo 周围留出空间 */
  will-change: filter;  /* 优化浏览器渲染性能 */
  transition: filter 300ms;  /* 滤镜效果过渡动画，持续 300 毫秒 */
}

/* logo 悬停效果 - 鼠标移到 logo 上时的样式变化 */
.logo:hover {
  /* drop-shadow 创建阴影效果 */
  filter: drop-shadow(0 0 2em #646cffaa);
}

/* .logo.vue 特定类 - Vue 官方 Logo 的悬停效果
   使用绿色 (#42b883) 作为品牌色 */
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
