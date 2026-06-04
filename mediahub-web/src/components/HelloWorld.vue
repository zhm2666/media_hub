<!--
/**
 * HelloWorld.vue - 示例组件
 *
 * 功能说明:
 * 这是一个 Vue 3 示例组件，展示 Vue 的基础用法。
 * 组件包含一个计数器功能，点击按钮可以增加计数。
 */
-->

<!-- ============================================================
     <script setup> 区域
     ============================================================ -->
<script setup lang="ts">
// 导入 Vue 的 ref 函数，用于创建响应式数据
import { ref } from 'vue'

// defineProps 是 Vue 3 的宏命令，用于定义组件的输入属性
// msg 属性可以从父组件传入
defineProps<{ msg: string }>()

// ref(0) 创建一个初始值为 0 的响应式变量
// 当数据变化时，视图会自动更新
const count = ref(0)
</script>

<!-- ============================================================
     <template> 区域
     ============================================================ -->
<template>
  <!-- 显示父组件传入的 msg 属性 -->
  <h1>{{ msg }}</h1>

  <!-- 卡片容器 -->
  <div class="card">
    <!--
      button 按钮
      type="button" - 防止表单提交
      @click="count++" - 点击事件，点击时计数器加 1
    -->
    <button type="button" @click="count++">count is {{ count }}</button>

    <!-- 提示文本 -->
    <p>
      Edit
      <code>components/HelloWorld.vue</code> to test HMR
    </p>
  </div>

  <!-- 文档链接 -->
  <p>
    Check out
    <a href="https://vuejs.org/guide/quick-start.html#local" target="_blank"
      >create-vue</a
    >, the official Vue + Vite starter
  </p>

  <p>
    Install
    <a href="https://github.com/vuejs/language-tools" target="_blank">Volar</a>
    in your IDE for a better DX
  </p>

  <!-- 文档提示 -->
  <p class="read-the-docs">Click on the Vite and Vue logos to learn more</p>
</template>

<!-- ============================================================
     <style scoped> 区域
     scoped 属性确保样式只影响当前组件
     ============================================================ -->
<style scoped>
/* 灰色文字样式，用于次要信息 */
.read-the-docs {
  color: #888;
}
</style>
