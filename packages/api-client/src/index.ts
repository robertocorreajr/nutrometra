export { api, ApiError, configureApiClient } from "./client"
export { ApiProvider } from "./provider"
export type { UUID, ErrorResponse, PaginatedResponse } from "./types/common"
export type { User, TenantMembership } from "./types/auth"
export type { Plan, Subscription, Entitlement } from "./types/billing"
export type {
  Patient,
  CreatePatientRequest,
  PatientProfile,
  PatientInvite,
  Gender,
} from "./types/patient"
