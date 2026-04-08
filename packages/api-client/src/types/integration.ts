export interface GoogleCalendarStatus {
  connected: boolean
  provider: string
  status: string
  synced_at?: string
}

export interface GoogleAuthorizeResponse {
  authorize_url: string
}
