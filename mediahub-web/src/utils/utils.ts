/**
 * utils.ts - 工具函数模块
 *
 * 模块说明:
 * 本模块提供通用的工具函数，包括 URL 参数获取、Cookie 操作、日期格式化等。
 * 这些函数可以在项目的任何地方复用。
 */

// ============================================================
// 函数 1: getUrlParameter - 获取 URL 参数
// ============================================================

/**
 * 从当前页面 URL 的查询字符串中获取指定参数的值
 *
 * @param name - 要获取的参数名称
 * @returns 返回参数的值，如果参数不存在则返回空字符串
 *
 * 示例:
 * URL: https://example.com?sys=ai&code=123
 * getUrlParameter('sys') => 'ai'
 */
export function getUrlParameter(name: string) {
  // 转义参数名中的特殊字符 [ 和 ]
  name = name.replace(/[[]/, '\\[').replace(/[\]]/, '\\]');

  // 创建正则表达式匹配 URL 参数
  // [\\?&] 匹配 ? 或 &，表示参数开始
  // name 是参数名
  // =([^&#]*) 捕获参数值（不含 # 和 &）
  var regex = new RegExp('[\\?&]' + name + '=([^&#]*)');

  // 执行正则匹配
  var results = regex.exec(window.location.search);

  // 如果没有匹配到，返回空字符串
  // 否则解码并返回参数值（处理 URL 编码中的 + 号为空格）
  return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
}

// ============================================================
// 函数 2: getCookie - 获取 Cookie 值
// ============================================================

/**
 * 从浏览器的 Cookie 中获取指定键的值
 *
 * @param key - 要获取的 Cookie 键名
 * @returns 返回 Cookie 的值，如果不存在则返回 null
 */
export function getCookie(key: string) {
  // 获取所有 Cookie 并分割成数组
  const cookies = document.cookie.split(';');

  // 遍历 Cookie 数组查找目标
  for (let i = 0; i < cookies.length; i++) {
    // 去除前后空格
    const cookie = cookies[i].trim();

    // 检查是否以 "key=" 开头
    if (cookie.startsWith(`${key}=`)) {
      // 返回 Cookie 值（去掉 "key=" 部分）
      return cookie.substring(key.length + 1);
    }
  }

  // 未找到返回 null
  return null;
}

// ============================================================
// 函数 3: getDateStr - 时间戳转日期字符串
// ============================================================

/**
 * 将 Unix 时间戳（秒）转换为格式化的日期字符串
 *
 * @param timestamp - Unix 时间戳（单位：秒）
 * @returns 格式化的日期字符串，格式：YYYY-MM-DD HH:mm:ss
 *
 * 示例:
 * getDateStr(1625123456) => "2021-07-01 12:30:56"
 */
export function getDateStr(timestamp: number) {
  // --------------------------------------------------------
  // 步骤 1: 将时间戳转换为 Date 对象
  // --------------------------------------------------------
  // JavaScript 的 Date 使用毫秒，而 Unix 时间戳是秒
  // 所以需要 * 1000 进行转换
  const date = new Date(timestamp * 1000);

  // --------------------------------------------------------
  // 步骤 2: 获取各个时间组成部分
  // --------------------------------------------------------
  const year = date.getFullYear();  // 年份（如 2021）

  // getMonth() 返回 0-11，需要 +1
  // ("0" + (month + 1)).slice(-2) 确保月份始终是两位数
  // 例如：1 月 => "01"，12 月 => "12"
  const month = ("0" + (date.getMonth() + 1)).slice(-2);

  // getDate() 返回日期（1-31）
  // 同样用 slice(-2) 确保两位数格式
  const day = ("0" + date.getDate()).slice(-2);

  // getHours() 返回小时（0-23）
  const hours = ("0" + date.getHours()).slice(-2);

  // getMinutes() 返回分钟（0-59）
  const minutes = ("0" + date.getMinutes()).slice(-2);

  // getSeconds() 返回秒数（0-59）
  const seconds = ("0" + date.getSeconds()).slice(-2);

  // --------------------------------------------------------
  // 步骤 3: 组装格式化字符串
  // --------------------------------------------------------
  let formattedTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  return formattedTime;
}

// ============================================================
// 函数 4: setCookie - 设置 Cookie
// ============================================================

/**
 * 设置一个 Cookie，包含过期时间
 *
 * @param name - Cookie 的名称
 * @param value - Cookie 的值
 * @param days - 过期天数（从当前日期开始计算）
 *
 * 示例:
 * setCookie('token', 'abc123', 7) // 设置一个 7 天后过期的 Cookie
 */
export function setCookie(name: string, value: string, days: number) {
  // 创建 Date 对象
  const date = new Date();

  // 计算过期时间：当前时间 + days 天
  // days * 24 * 60 * 60 * 1000 转换为毫秒
  date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));

  // 构建过期时间字符串（UTC 格式）
  const expires = "expires=" + date.toUTCString();

  // 设置 Cookie
  // path=/ 表示整个网站都可以访问这个 Cookie
  document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

// ============================================================
// 函数 5: deleteCookie - 删除 Cookie
// ============================================================

/**
 * 删除指定名称的 Cookie
 *
 * 原理：通过将过期时间设置为过去的日期，浏览器会自动删除该 Cookie
 *
 * @param name - 要删除的 Cookie 名称
 *
 * 示例:
 * deleteCookie('token') // 删除名为 token 的 Cookie
 */
export function deleteCookie(name: string) {
  // 设置 Cookie 的过期时间为 1970 年 1 月 1 日
  // 浏览器检测到过期时间已过，会自动删除该 Cookie
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`;
}
