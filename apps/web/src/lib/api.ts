import axios from 'axios'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL,
})

function getCorrelationId() {
  return `web-${crypto.randomUUID()}`
}

api.interceptors.request.use((config) => {
  config.headers['X-Correlation-ID'] = getCorrelationId()

  const debugUserId = localStorage.getItem('debug_user_id')
  const debugUserRole = localStorage.getItem('debug_user_role')

  if (debugUserId) {
    config.headers['X-Debug-User-ID'] = debugUserId
  }

  if (debugUserRole) {
    config.headers['X-Debug-User-Role'] = debugUserRole
  }

  return config
})