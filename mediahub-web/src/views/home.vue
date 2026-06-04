<!--
/**
 * home.vue - 首页组件
 *
 * 功能说明:
 * 这是应用的主页面，包含以下功能：
 * 1. 顶部导航菜单
 * 2. 轮播图展示区
 * 3. 图片网格展示区（两行）
 * 4. 文件上传功能
 *
 * 技术栈:
 * - Vue 3 Composition API
 * - Element Plus 组件库
 */
-->

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

// logo - Logo 组件
import logo from "./components/logo.vue";

// upload - 文件上传组件
import upload from "./components/upload.vue"

// onBeforeMount - 生命周期钩子，组件挂载前执行
// ref - 创建响应式数据
import { onBeforeMount, ref } from 'vue'

// home - 获取首页数据的 API 函数
import { home } from '../api/api.ts'

// menu1 - 顶部导航菜单组件
import menu1 from "./menu.vue"

// ============================================================
// 2. 定义响应式数据
// ============================================================

/**
 * data - 页面数据对象
 *
 * 包含：
 * - banners - 轮播图 URL 数组
 * - imgs1 - 第一行图片 URL 数组
 * - imgs2 - 第二行图片 URL 数组
 */
let data = {
  // 轮播图，初始三个空字符串（用于骨架屏/占位）
  banners: ref(["", "", ""]),

  // 第一行图片，初始五个空字符串
  imgs1: ref(["", "", "", "", ""]),

  // 第二行图片，初始五个空字符串
  imgs2: ref(["", "", "", "", ""]),
};

// ============================================================
// 3. 定义数据类型
// ============================================================

/**
 * homeRes - 首页响应数据类型
 *
 * 使用 ES6 Class 定义，匹配后端返回的数据结构
 */
class homeRes {
  // 轮播图 URL 数组
  banners: Array<string>;

  // 第一行图片 URL 数组
  images1: Array<string>;

  // 第二行图片 URL 数组
  images2: Array<string>;

  // 构造函数
  constructor(banners: Array<string>, images1: Array<string>, images2: Array<string>) {
    this.banners = banners;
    this.images1 = images1;
    this.images2 = images2;
  }
}

// ============================================================
// 4. 组件挂载前生命周期钩子
// ============================================================

/**
 * onBeforeMount - 组件挂载到 DOM 前执行
 *
 * 用于获取首页数据
 */
onBeforeMount(() => {
  console.log("on before mount");  // 调试日志

  // 调用首页 API
  home<homeRes>().then(function (res) {
    console.log(res.data);  // 调试：打印响应数据

    // 更新轮播图数据
    // res.data?.banners 使用可选链，如果 data 或 banners 为 undefined 则使用默认值
    data.banners.value = res.data?.banners || ["", "", ""];

    // 更新第一行图片
    data.imgs1.value = res.data?.images1 || ["", "", "", "", ""];

    // 更新第二行图片
    data.imgs2.value = res.data?.images2 || ["", "", "", "", ""];

  }).catch(function (res) {
    // 错误处理
    console.log(res);
  });
});

// ============================================================
// 5. 样式计算函数
// ============================================================

/**
 * getStyle - 计算图片容器的 CSS 样式
 *
 * 用于动态设置图片网格中每个格子的位置和尺寸
 *
 * @param index - 图片索引（0-4）
 * @returns CSS 样式字符串
 */
function getStyle(index: number) {
  // 基础样式
  let style = "position:absolute;width:14rem;height:10rem;top:0;background-color: #d3dce6;";

  // 计算 left 位置：每个格子间隔 15.25rem
  // index=0: left=0
  // index=1: left=15.25rem
  // index=2: left=30.5rem
  // 以此类推...
  style += " left:" + 15.25 * index + 'rem;';

  return style;
}

/**
 * getItemStyle - 计算图片项的 CSS 样式
 *
 * 用于设置背景图片的显示方式
 *
 * @param url - 图片 URL
 * @returns CSS 样式字符串
 */
function getItemStyle(url: string) {
  // 基础尺寸
  let style = "width:100%;height:100%;";

  // background-image - 背景图片 URL
  style += "background-image:url(" + url + ");";

  // background-position - 图片在容器中的位置（居中）
  style += "background-position:center center;";

  // background-repeat - 不重复平铺
  style += "background-repeat:no-repeat;";

  // background-size - 覆盖整个容器（保持比例）
  style += "background-size:cover;";

  return style;
}
</script>

<!-- ============================================================
     <template> 区域
     ============================================================ -->
<template>
  <!-- 顶部导航菜单 -->
  <menu1></menu1>

  <!-- 页面主体容器 -->
  <!-- position:relative 作为定位参考 -->
  <!-- width/height 设置页面高度 -->
  <div style="position:relative;width:100%;height:57.1875rem;">

    <!-- --------------------------------------------------------
         轮播图区域
         -------------------------------------------------------- -->
    <!-- el-carousel - Element Plus 轮播图组件 -->
    <!-- width/height 设置轮播图尺寸 -->
    <!-- :interval="5000" 自动切换间隔 5 秒 -->
    <!-- arrow="always" 始终显示左右箭头 -->
    <el-carousel width="100%" height="25.9375rem" class="banner" :interval="5000" arrow="always">
      <!-- v-for 循环渲染轮播项 -->
      <!-- v-for="item in data.banners.value" 遍历轮播图数组 -->
      <!-- :key="item" 为每个元素提供唯一标识 -->
      <el-carousel-item v-for="item in data.banners.value" :key="item">
        <!-- 使用背景图样式显示轮播图 -->
        <!-- :style 绑定动态计算的样式 -->
        <div :style="getItemStyle(item)"></div>
      </el-carousel-item>
    </el-carousel>

    <!-- --------------------------------------------------------
         Logo 和上传区域（覆盖在轮播图上方）
         -------------------------------------------------------- -->
    <!-- banner_upper 容器包含 Logo 和上传按钮 -->
    <!-- position:absolute 使其浮在轮播图上方 -->
    <div class="banner_upper">
      <!-- Logo 组件 -->
      <!-- position:absolute;top:0 定位到容器顶部 -->
      <logo style="position: absolute;top:0;"></logo>

      <!-- 上传组件 -->
      <!-- position:absolute;top:7.663rem 定位到 Logo 下方 -->
      <upload style="position: absolute;top:7.663rem"></upload>
    </div>

    <!-- --------------------------------------------------------
         图片展示区域
         -------------------------------------------------------- -->
    <!-- 外层容器：flex 布局使内容居中 -->
    <div style="width: 100%;height:21.875rem; position:relative;top:5rem;display:flex;justify-content: center;">
      <!-- 内层容器：固定宽度 -->
      <div style="position: relative;width:75rem;height:21.875rem;">

        <!-- 第一行图片 -->
        <div style="position: relative;width:100%;height:10rem">
          <!-- v-for 循环渲染 5 个图片 -->
          <div v-for="(item,index) in data.imgs1.value" :style="getStyle(index)">
            <!-- 每个图片使用背景图样式 -->
            <div :style="getItemStyle(item)"></div>
          </div>
        </div>

        <!-- 第二行图片 -->
        <!-- top:1.875rem 与第一行保持间距 -->
        <div style="position: relative;width:100%;height:10rem;top:1.875rem">
          <div v-for="(item,index) in data.imgs2.value" :style="getStyle(index)">
            <div :style="getItemStyle(item)"></div>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<!-- ============================================================
     <style scoped> 区域
     ============================================================ -->
<style scoped>
/* --------------------------------------------------------
   轮播图项样式
   -------------------------------------------------------- */

/* el-carousel__item h3 - 轮播图项内的标题样式 */
.el-carousel__item h3 {
  color: #475669;        /* 深蓝灰色文字 */
  opacity: 0.75;         /* 75% 透明度 */
  line-height: 580px;   /* 行高与轮播图高度一致 */
  margin: 0;             /* 去除默认外边距 */
  text-align: center;    /* 文字居中 */
}

/* 偶数项背景色 */
.el-carousel__item:nth-child(2n) {
  background-color: #99a9bf;  /* 蓝灰色 */
}

/* 奇数项背景色 */
.el-carousel__item:nth-child(2n + 1) {
  background-color: #d3dce6;  /* 浅蓝灰色 */
}

/* --------------------------------------------------------
   .banner 轮播图容器样式
   -------------------------------------------------------- */
.banner {
  width: 100%;              /* 宽度占满 */
  min-height: 25.9375rem;   /* 最小高度 */
  position: relative;        /* 相对定位 */
}

/* --------------------------------------------------------
   .banner_upper Logo/上传区域容器样式
   -------------------------------------------------------- */
.banner_upper {
  display: flex;                     /* 弹性盒子布局 */
  width: 46.9%;                      /* 宽度 46.9% */
  height: 11.5rem;                   /* 高度 */
  z-index: 1;                        /* 层级，高于轮播图 */
  top: 6rem;                         /* 距离顶部 6rem */
  left: 26.8%;                       /* 距离左侧 26.8% */
  position: absolute;                 /* 绝对定位 */
  text-align: center;                /* 文本居中 */
  justify-content: center;           /* 水平居中 */
}
</style>
