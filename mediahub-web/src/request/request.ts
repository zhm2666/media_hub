/**
 * request.ts - HTTP 请求封装模块
 *
 * 模块说明:
 * 本模块对 axios 进行了二次封装，提供更简洁的 HTTP 请求方法。
 * 支持 GET、POST 请求，并统一了请求参数和响应格式。
 */

// ============================================================
// 1. 导入依赖
// ============================================================

// 导入配置好的 axios 实例
import request from "./axios"

// 导入 axios 的类型定义
import type {
  AxiosProgressEvent,  // 进度事件类型
  AxiosResponse,       // 响应对象类型
  GenericAbortSignal,   // 请求取消信号类型
  AxiosError            // 错误类型
} from "axios";

// ============================================================
// 2. 定义类型接口
// ============================================================

/**
 * HttpParams - HTTP 请求参数接口
 *
 * 定义所有 HTTP 请求可用的配置选项。
 */
export interface HttpParams {
  /** 请求的 URL 地址 */
  url: string;

  /** 请求数据（GET 时为查询参数，POST 时为请求体） */
  data?: any;

  /** 请求方法，默认为 'GET' 或 'POST' */
  method?: string;

  /** 自定义请求头 */
  headers?: any;

  /** 下载进度回调 */
  onDownloadProgress?: (progressEvent: AxiosProgressEvent) => void;

  /** 上传进度回调 */
  onUploadProgress?: (progressEvent: AxiosProgressEvent) => void;

  /** 请求取消信号（用于 AbortController） */
  signal?: GenericAbortSignal;

  /** 请求发送前的回调 */
  beforeRequest?: () => void;

  /** 请求完成后的回调（无论成功或失败） */
  afterRequest?: () => void;
}

/**
 * Response - 统一封装的响应数据结构
 *
 * @template T - data 的类型
 */
export interface Response<T = any> {
  /** 响应数据 */
  data?: T;

  /** HTTP 状态码 */
  status: number;

  /** 响应消息（用于携带错误信息） */
  message?: string;
}

// ============================================================
// 3. 核心 HTTP 请求函数
// ============================================================

/**
 * http - 通用 HTTP 请求函数
 *
 * 处理 GET 和 POST 请求的核心逻辑。
 *
 * @template T - 响应数据的类型
 */
function http<T = any>({
  url,
  data,
  method,
  headers,
  onDownloadProgress,
  onUploadProgress,
  signal,
  beforeRequest,
  afterRequest
}: HttpParams): Promise<Response<T>> {
  // --------------------------------------------------------
  // 成功处理函数
  // --------------------------------------------------------
  const successHandler = (res: AxiosResponse<T>) => {
    // 返回标准化格式：{ data: 服务器数据, status: HTTP状态码 }
    return Promise.resolve({ data: res.data, status: res.status });
  };

  // --------------------------------------------------------
  // 失败处理函数
  // --------------------------------------------------------
  const failHandler = (err: Error) => {
    const axiosErr = err as AxiosError;
    // 返回被拒绝的 Promise，包含状态码和错误消息
    return Promise.reject({
      status: axiosErr.response?.status,
      message: axiosErr.message
    });
  };

  // --------------------------------------------------------
  // 执行请求前回调
  // --------------------------------------------------------
  beforeRequest?.();

  // --------------------------------------------------------
  // 设置默认请求方法
  // --------------------------------------------------------
  method = method || 'GET';

  // --------------------------------------------------------
  // 处理请求数据
  // --------------------------------------------------------
  // 如果 data 是函数，调用它获取数据
  // 否则使用 data 或空对象
  const params = Object.assign(
    typeof data === 'function' ? data() : data ?? {},
    {}
  );

  // --------------------------------------------------------
  // 根据请求方法发送请求
  // --------------------------------------------------------
  return method === "GET"
    // GET 请求
    ? request.get(url, { params, headers, signal, onDownloadProgress })
        .then(successHandler, failHandler)
        .finally(afterRequest)
    // POST 请求
    : request.post(url, params, { headers, signal, onDownloadProgress, onUploadProgress })
        .then(successHandler, failHandler)
        .finally(afterRequest);
}

// ============================================================
// 4. GET 请求封装
// ============================================================

/**
 * get - GET 请求封装函数
 *
 * @template T - 响应数据的类型
 */
export function get<T = any>(
  { url, data, method = "GET", headers, onDownloadProgress, signal, beforeRequest, afterRequest }: HttpParams
): Promise<Response<T>> {
  return http<T>({
    url, method, data, headers, onDownloadProgress, signal, beforeRequest, afterRequest
  });
}

// ============================================================
// 5. POST 请求封装
// ============================================================

/**
 * post - POST 请求封装函数
 *
 * @template T - 响应数据的类型
 */
export function post<T = any>(
  { url, data, method = 'POST', headers, onDownloadProgress, onUploadProgress, signal, beforeRequest, afterRequest }: HttpParams
): Promise<Response<T>> {
  return http<T>({
    url, method, data, headers, onDownloadProgress, onUploadProgress, signal, beforeRequest, afterRequest
  });
}

// ============================================================
// 6. 默认导出
// ============================================================

// 默认导出 post 函数
export default post;
