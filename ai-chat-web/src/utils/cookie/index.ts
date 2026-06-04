export function getCookieValue(key: string) {
  const cookies = document.cookie.split(';')
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim()
    if (cookie.startsWith(`${key}=`))
      return cookie.substring(key.length + 1)
  }
  return null
}

// getCookieByKey 是 getCookieValue 的别名，用于兼容 axios.ts 中的引用
export function getCookieByKey(key: string) {
  return getCookieValue(key)
}

export function deleteCookieByKey(key: string) {
  document.cookie = `${key}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
}
