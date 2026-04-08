export type UUID = string

export interface ErrorResponse {
  code: string
  message: string
  request_id: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  limit: number
  offset: number
}
