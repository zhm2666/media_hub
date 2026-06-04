/**
 * vite-env.d.ts - Vite 环境类型声明文件
 *
 * 文件说明:
 * 这是 TypeScript 的类型声明文件（.d.ts = declaration types）
 * 主要用于告诉 TypeScript 编译器 Vite 提供的一些全局变量和类型
 *
 * 为什么要这个文件:
 * Vite 在构建时会在项目中注入一些全局变量（如 import.meta.env），
 * 但 TypeScript 默认不认识这些变量，会报错。通过声明文件，
 * 我们告诉 TypeScript 这些变量是存在的，避免类型错误。
 */

/**
 * /// <reference types="vite/client" />
 *
 * 这是一个三斜线指令（Triple-slash directive）
 * 作用是引用外部的类型声明文件
 *
 * "vite/client" 包含以下类型声明:
 * - import.meta.env - Vite 的环境变量对象
 * - import.meta.env.PROD - 是否生产环境
 * - import.meta.env.DEV - 是否开发环境
 * - import.meta.env.BASE_URL - 应用的基础 URL
 * - import.meta.glob - 动态导入文件的 glob 模式
 * - Vite HMR（热模块替换）相关的类型
 *
 * 常见用法:
 * 在 .env 文件中定义变量如 VITE_API_BASE_URL
 * 然后在代码中使用 import.meta.env.VITE_API_BASE_URL 访问
 */
/// <reference types="vite/client" />
