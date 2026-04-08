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
export type { Professional } from "./types/professional"
export type {
  Anamnesis,
  AnamnesisStatus,
  AnamnesisRequest,
  ProgressNote,
  ProgressNoteRequest,
  AttachmentCategory,
  ClinicalAttachment,
  AttachmentRequest,
} from "./types/clinical"
export type {
  MeasurementSource,
  BodyMeasurement,
  CreateMeasurementRequest,
} from "./types/bioimpedance"
export type {
  SuggestionType,
  SuggestionStatus,
  AISuggestion,
  CreateSuggestionRequest,
} from "./types/ai"
export type {
  ServiceMode,
  AppointmentStatus,
  AppointmentSource,
  AvailabilityRule,
  CreateAvailabilityRuleRequest,
  ScheduleBlock,
  CreateBlockRequest,
  TimeSlot,
  Appointment,
  CreateAppointmentRequest,
  UpdateStatusRequest,
  RescheduleRequest,
} from "./types/scheduling"
