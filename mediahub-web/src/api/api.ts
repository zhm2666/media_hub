/**
 * api.ts - API 接口封装模块
 *
 * 模块说明:
 * 本模块封装了项目所需的所有后端 API 接口。
 * 提供文件上传和首页数据获取功能。
 */

// ============================================================
// 1. 导入 HTTP 请求工具
// ============================================================
// 从 request.ts 导入封装好的 get 和 post 方法
import { get, post } from '../request/request'

// ============================================================
// 2. 文件上传接口
// ============================================================

/**
 * uploadFile - 上传文件
 *
 * 功能说明:
 * 将文件上传到服务器，返回文件的访问地址。
 *
 * @template T - 返回数据的类型
 * @param {Object} params - 请求参数
 * @param {FormData} params.formData - 包含文件的 FormData 对象
 *   FormData 是浏览器提供的 API，用于构造 multipart/form-data 请求
 * @returns {Promise} 返回 Promise 对象，包含文件 URL
 *
 * API 详情:
 * - 请求方法: POST
 * - 请求路径: /v1/file/upload
 * - Content-Type: multipart/form-data（自动设置）
 *
 * 使用示例:
 * ```typescript
 * const formData = new FormData();
 * formData.append('file', fileInput.files[0]);
 * uploadFile<{ url: string }>({ formData })
 *   .then(res => console.log(res.data.url));
 * ```
 */
export function uploadFile<T = any>(params: { formData: FormData }) {
  // 定义 API 路径
  const path = '/v1/file/upload';

  // 调用 post 方法发送请求
  // formData 会作为请求体自动以 multipart/form-data 格式发送
  return post<T>({
    url: path,
    data: params.formData,
  });
}

// ============================================================
// 3. 首页数据接口
// ============================================================

/**
 * home - 获取首页数据
 *
 * 功能说明:
 * 获取首页所需的轮播图和图片数据。
 *
 * @template T - 返回数据的类型
 * @returns {Promise} 返回 Promise 对象，包含首页数据
 *
 * API 详情:
 * - 请求方法: GET
 * - 请求路径: /v1/home
 *
 * 返回数据格式:
 * {
 *   banners: ["url1", "url2", "url3"],      // 轮播图 URL 数组
 *   images1: ["url1", "url2", ...],         // 第一行图片 URL 数组
 *   images2: ["url1", "url2", ...]          // 第二行图片 URL 数组
 * }
 *
 * 使用示例:
 * ```typescript
 * home<HomeData>().then(res => {
 *   this.banners = res.data.banners;
 * });
 * ```
 */
export function home<T = any>() {
  // 定义 API 路径
  const path = '/v1/home';

  // 调用 get 方法发送请求
  return get<T>({ url: path });
}
