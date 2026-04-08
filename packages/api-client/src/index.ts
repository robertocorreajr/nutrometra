export { api, ApiError, configureApiClient } from "./client"
export { ApiProvider } from "./provider"
export type { UUID, ErrorResponse, PaginatedResponse } from "./types/common"
export type { User, TenantMembership } from "./types/auth"
export type {
  Plan, Subscription, Entitlement,
  InvoiceStatus, Invoice, PaymentStatus, Payment, CheckoutRequest, ChangePlanRequest,
} from "./types/billing"
export type {
  Patient,
  CreatePatientRequest,
  PatientProfile,
  PatientInvite,
  Gender,
} from "./types/patient"
export type {
  Professional,
  UpdateProfessionalRequest,
  ProfessionalAddress,
  CreateAddressRequest,
  ServiceModeType,
  ProfessionalServiceMode,
  SetServiceModeRequest,
} from "./types/professional"
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
export type {
  FoodSource,
  NutritionFacts,
  HouseholdMeasure,
  FoodItem,
  CreateFoodItemRequest,
} from "./types/catalog"
export type {
  DietStatus,
  DietSubstitution,
  DietMealItem,
  DietMeal,
  Diet,
  CreateDietRequest,
  UpdateDietRequest,
  CreateMealRequest,
  CreateMealItemRequest,
  CreateSubstitutionRequest,
} from "./types/diet"
export type {
  DocumentType,
  DocumentStatus,
  ClinicalDocument,
  CreateDocumentRequest,
  UpdateDocumentRequest,
} from "./types/document"
export type {
  ExportStatus,
  ExportType,
  ExportedFile,
  RequestExportRequest,
} from "./types/export"
export type { GoogleCalendarStatus, GoogleAuthorizeResponse } from "./types/integration"
export type { Role, AssignRoleRequest } from "./types/rbac"
