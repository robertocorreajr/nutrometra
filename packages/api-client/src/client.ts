import type { ErrorResponse } from "./types/common"

export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    public requestId: string,
  ) {
    super(`API Error: ${code}`)
    this.name = "ApiError"
  }
}

interface ApiClientConfig {
  baseUrl: string
  getAccessToken: () => Promise<string | null>
  getTenantId: () => string | null
}

let clientConfig: ApiClientConfig | null = null

export function configureApiClient(config: ApiClientConfig) {
  clientConfig = config
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  if (!clientConfig) {
    throw new Error("API client not configured. Call configureApiClient() first.")
  }

  const headers = new Headers(options.headers)
  headers.set("Content-Type", "application/json")

  const token = await clientConfig.getAccessToken()
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const tenantId = clientConfig.getTenantId()
  if (tenantId) {
    headers.set("X-Tenant-ID", tenantId)
  }

  headers.set("X-Request-ID", crypto.randomUUID())

  const res = await fetch(`${clientConfig.baseUrl}${path}`, {
    ...options,
    headers,
  })

  if (!res.ok) {
    const body: ErrorResponse = await res.json().catch(() => ({
      code: "unknown",
      message: res.statusText,
      request_id: "",
    }))
    throw new ApiError(res.status, body.code, body.request_id)
  }

  if (res.status === 204) return undefined as T

  return res.json()
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, { method: "POST", body: body ? JSON.stringify(body) : undefined }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PATCH", body: JSON.stringify(body) }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),

  download: async (path: string, filename: string) => {
    if (!clientConfig) throw new Error("API client not configured. Call configureApiClient() first.")
    const headers = new Headers()
    const token = await clientConfig.getAccessToken()
    if (token) headers.set("Authorization", `Bearer ${token}`)
    const tenantId = clientConfig.getTenantId()
    if (tenantId) headers.set("X-Tenant-ID", tenantId)
    headers.set("X-Request-ID", crypto.randomUUID())
    const res = await fetch(`${clientConfig.baseUrl}${path}`, { headers })
    if (!res.ok) throw new ApiError(res.status, "download_failed", "")
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)
  },
}
