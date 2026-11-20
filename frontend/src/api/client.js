import axios from 'axios'

/**
 * Resolves the API base URL dynamically:
 * - If running in browser on production domain (not localhost):
 *   Always use relative path '/api/v1' so requests route securely via Nginx reverse proxy (HTTPS).
 *   This prevents hardcoded 'http://localhost:8088' from breaking in live environments (Mixed Content block).
 * - If running in local development:
 *   Use VITE_API_URL if defined, or '/api/v1' (which is proxied by Vite dev server to port 8088).
 */
export function getApiBaseUrl() {
  const envUrl = (import.meta.env?.VITE_API_URL || '').trim()

  if (typeof window !== 'undefined' && window.location?.hostname) {
    const host = window.location.hostname.toLowerCase()
    const isLocal = host === 'localhost' || host === '127.0.0.1' || host === '::1' || host.endsWith('.local') || host.endsWith('.test')

    // Di environment server live / production:
    // Jangan pernah izinkan localhost:8088 karena browser akan memblokir Mixed Content (HTTPS -> HTTP)
    if (!isLocal) {
      if (envUrl && !envUrl.includes('localhost') && !envUrl.includes('127.0.0.1')) {
        return envUrl.replace(/\/$/, '')
      }
      return '/api/v1'
    }
  }

  return envUrl || '/api/v1'
}

const apiClient = axios.create({
  baseURL: getApiBaseUrl(),
  headers: {
    'Content-Type': 'application/json',
  },
})

export default apiClient

