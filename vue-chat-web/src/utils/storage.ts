export const storage = {
  getToken: (): string | null => {
    return localStorage.getItem('chat_token')
  },

  setToken: (token: string): void => {
    localStorage.setItem('chat_token', token)
  },

  removeToken: (): void => {
    localStorage.removeItem('chat_token')
  },

  getUser: (): any => {
    const user = localStorage.getItem('chat_user')
    return user ? JSON.parse(user) : null
  },

  setUser: (user: any): void => {
    localStorage.setItem('chat_user', JSON.stringify(user))
  },

  removeUser: (): void => {
    localStorage.removeItem('chat_user')
  },

  clear: (): void => {
    localStorage.removeItem('chat_token')
    localStorage.removeItem('chat_user')
  },
}
