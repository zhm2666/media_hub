/**
 * axios.ts - Axios 实例配置与拦截器
 *
 * 模块说明:
 * 本模块创建并配置了一个全局的 Axios 实例。
 * 包含请求拦截器（添加 Token）和响应拦截器（统一错误处理）。
 */

// ============================================================
// 1. 导入依赖
// ============================================================

// axios - HTTP 请求库
// AxiosError - axios 的错误类型
import axios, { AxiosError } from 'axios'

// ElMessage - Element Plus 的消息提示组件
import { ElMessage } from 'element-plus'

// getCookie - 从 utils.ts 导入的 Cookie 获取函数
import { getCookie } from "../utils/utils.ts";

// ============================================================
// 2. 创建 Axios 实例
// ============================================================

/**
 * axios.create() - 创建 Axios 实例
 *
 * baseURL - API 的基础地址
 * 来自 Vite 环境变量 VITE_API_BASE_URL
 */
const service = axios.create({
  // baseURL - API 基础地址
  // 例如: http://localhost:8080/api
  baseURL: import.meta.env.VITE_API_BASE_URL
});

// ============================================================
// 3. 请求拦截器 (Request Interceptor)
// ============================================================

/**
 * 请求拦截器在请求发送之前执行
 * 用于统一添加认证 Token 等
 */
service.interceptors.request.use(
  /**
   * @param config - Axios 的请求配置对象
   * @returns 返回修改后的配置
   */
  function (config) {
    // --------------------------------------------------------
    // 从 Cookie 中获取访问令牌
    // --------------------------------------------------------
    // sso_0voice_access_token 是后端设置的认证 Token
    // 存储在 Cookie 中，前端可以方便地获取
    let access_token = getCookie("sso_0voice_access_token");

    // --------------------------------------------------------
    // 如果 Token 存在，添加到请求头
    // --------------------------------------------------------
    if (access_token) {
      // Authorization 请求头使用 Bearer 方案
      // 格式: Authorization: Bearer <token>
      // 这是 OAuth 2.0 的标准授权方式
      config.headers.Authorization = "Bearer " + access_token;
    }

    // 返回配置，请求会使用修改后的配置发送
    return config;
  },

  /**
   * 请求发送失败时的处理
   * @param error - 错误对象
   */
  function (error) {
    // 打印错误
    console.error('请求发送失败:', error);

    // 返回被拒绝的 Promise
    return Promise.reject(error);
  }
);

// ============================================================
// 4. 响应拦截器 (Response Interceptor)
// ============================================================

/**
 * 响应拦截器在收到响应后执行
 * 用于统一处理错误
 */
service.interceptors.response.use(
  /**
   * 响应成功处理
   * @param response - Axios 响应对象
   */
  function (response) {
    // 2xx 状态码直接返回响应
    return response;
  },

  /**
   * 响应失败处理
   * @param error - 错误对象
   */
  function (error) {
    // 转换为 AxiosError 类型
    const axiosErr = error as AxiosError;

    // ============================================================
    // 根据 HTTP 状态码显示不同的错误提示
    // ============================================================

    // --------------------------------------------------------
    // 没有收到服务器响应（网络错误等）
    // --------------------------------------------------------
    if (!axiosErr.response?.status) {
      ElMessage({
        showClose: true,
        message: axiosErr.message,
        type: 'error',
      });
    }

    // --------------------------------------------------------
    // HTTP 500 - 服务器内部错误
    // --------------------------------------------------------
    else if (axiosErr.response?.status == 500) {
      ElMessage({
        showClose: true,
        message: "服务器内部错误",
        type: 'error',
      });
    }

    // --------------------------------------------------------
    // HTTP 504 - 网关超时
    // --------------------------------------------------------
    else if (axiosErr.response?.status == 504) {
      ElMessage({
        showClose: true,
        message: "网关超时",
        type: 'error',
      });
    }

    // --------------------------------------------------------
    // HTTP 413 - 请求体过大
    // --------------------------------------------------------
    else if (axiosErr.response?.status == 413) {
      ElMessage({
        showClose: true,
        message: "仅支持上传20M以内的图片",
        type: 'error',
      });
    }

    // --------------------------------------------------------
    // 打印调试信息
    // --------------------------------------------------------
    console.log(axiosErr.message);
    console.log(axiosErr.response?.status);

    // 返回被拒绝的 Promise
    return Promise.reject(error);
  }
);

// ============================================================
// 5. 导出 Axios 实例
// ============================================================
export default service;
