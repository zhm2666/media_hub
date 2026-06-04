<!--
/**
 * logo.vue - Logo 组件
 *
 * 功能说明:
 * 显示 MEDIAHUB 的品牌 Logo。
 * 使用 CSS 边框和定位创建有层次感的视觉效果。
 *
 * 视觉效果:
 * 多层嵌套的边框盒子，最内层显示文字 "MEDIAHUB"。
 */

<!-- ============================================================
     <template> 区域
     ============================================================ -->
<template>
  <!-- Logo 外层容器 -->
  <div>
    <!--
      外层 div - 定义 Logo 的总尺寸
      width: 17.8125rem (约 285px)
      height: 3.375rem (约 54px)
    -->
    <div style="width: 17.8125rem;height:3.375rem">

      <!-- --------------------------------------------------------
           第一层边框
           -------------------------------------------------------- -->
      <!--
        使用绝对定位覆盖在外层容器上
        border: 0.125rem solid white - 2px 白色实线边框
        border-style: solid - 实线边框样式
      -->
      <div
        style="position: absolute;top:0;left:0;width:17.125rem;height:2.75rem;border:0.125rem;border-style: solid;border-color: white;">

        <!-- --------------------------------------------------------
             第二层边框（嵌套在第一层内）
             -------------------------------------------------------- -->
        <!--
          向内偏移 0.25rem
          宽度/高度各减少 0.5rem（左右各 0.25rem）
          同样白色边框，形成双层边框效果
        -->
        <div
          style="position: absolute;top:0.25rem;left:0.25rem;width:17.125rem;height:2.75rem;border:0.125rem;border-style: solid;border-color: white;">

          <!-- --------------------------------------------------------
               文字区域
               -------------------------------------------------------- -->
          <!--
            文字容器
            position:absolute 绝对定位
            top/left 调整位置到边框内部
          -->
          <div
            style="position:absolute;top:0.5rem;left:0.875rem; height: 1.5rem;width:14.625rem;font-size: 1.5rem;font-weight:bolder;letter-spacing: 0.875rem;color:white;line-height: 1.5rem;">
            <!--
              文字样式说明:
              - font-size: 1.5rem - 字体大小 24px
              - font-weight:bolder - 比 bold 更粗
              - letter-spacing: 0.875rem - 字符间距 14px
              - color:white - 白色文字
              - line-height: 1.5rem - 行高 24px
            -->

            <!-- MEDIEHUB 文字 -->
            MEDIAHUB

          </div>
        </div>
      </div>
    </div>
  </div>
</template>
