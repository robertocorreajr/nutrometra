# Frontend Shared Infrastructure — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Set up the Turborepo + pnpm monorepo with 3 shared packages (@nutrometra/ui, @nutrometra/api-client, @nutrometra/auth) and wire all 3 Next.js apps to use them, resulting in working dev server with auth-protected layout shell on each app.

**Architecture:** Turborepo orchestrates 3 Next.js apps (`apps/`) + 3 internal TypeScript packages (`packages/`). Each package is a plain TS library with `src/` entrypoint. Apps import packages via pnpm workspace protocol. Tailwind preset shared from `@nutrometra/ui`.

**Tech Stack:** Next.js 14, React 18, TypeScript 5, Turborepo, pnpm, Tailwind CSS, shadcn/ui, TanStack Query, React Hook Form, Zod, next-auth (Zitadel OIDC provider).

**Spec:** `docs/superpowers/specs/2026-04-07-frontend-shared-infra-design.md`

---

## File Structure

### New files (root)

| File | Responsibility |
|------|---------------|
| `package.json` | Root workspace — scripts, devDependencies for tooling |
| `pnpm-workspace.yaml` | Declare workspace packages |
| `turbo.json` | Task pipelines (build, dev, lint) |
| `tsconfig.base.json` | Shared TS compiler options |
| `.npmrc` | pnpm settings (strict-peer-deps, etc) |
| `.node-version` | Pin Node 20 LTS |

### New files (packages/ui)

| File | Responsibility |
|------|---------------|
| `packages/ui/package.json` | Package manifest |
| `packages/ui/tsconfig.json` | TS config extending base |
| `packages/ui/src/index.ts` | Package entrypoint — re-exports |
| `packages/ui/src/tailwind/preset.ts` | Shared Tailwind preset (colors, fonts, breakpoints) |
| `packages/ui/src/lib/utils.ts` | cn() utility for shadcn/ui class merging |
| `packages/ui/src/components/button.tsx` | shadcn Button |
| `packages/ui/src/components/input.tsx` | shadcn Input |
| `packages/ui/src/components/card.tsx` | shadcn Card |
| `packages/ui/src/components/toast.tsx` | shadcn Toast + Toaster |
| `packages/ui/src/components/skeleton.tsx` | shadcn Skeleton |
| `packages/ui/src/layouts/app-shell.tsx` | Sidebar + header + content responsive layout |
| `packages/ui/src/layouts/page-header.tsx` | Title + breadcrumbs + actions |
| `packages/ui/src/layouts/empty-state.tsx` | Icon + message + CTA |
| `packages/ui/src/layouts/loading-state.tsx` | Skeleton-based loading |
| `packages/ui/src/layouts/error-state.tsx` | Error message + retry |

### New files (packages/api-client)

| File | Responsibility |
|------|---------------|
| `packages/api-client/package.json` | Package manifest |
| `packages/api-client/tsconfig.json` | TS config extending base |
| `packages/api-client/src/index.ts` | Re-exports |
| `packages/api-client/src/client.ts` | Fetch wrapper — auth headers, tenant header, error handling |
| `packages/api-client/src/provider.tsx` | React context provider — QueryClientProvider + ApiClientProvider |
| `packages/api-client/src/types/common.ts` | ErrorResponse, PaginatedResponse, UUID |
| `packages/api-client/src/types/auth.ts` | User, Session types |
| `packages/api-client/src/types/billing.ts` | Plan, Subscription, Entitlement |
| `packages/api-client/src/types/index.ts` | Re-exports all types |
| `packages/api-client/src/hooks/billing.ts` | usePlans, useSubscription, useEntitlements |
| `packages/api-client/src/hooks/index.ts` | Re-exports all hooks |

### New files (packages/auth)

| File | Responsibility |
|------|---------------|
| `packages/auth/package.json` | Package manifest |
| `packages/auth/tsconfig.json` | TS config extending base |
| `packages/auth/src/index.ts` | Re-exports |
| `packages/auth/src/config.ts` | NextAuth config — Zitadel OIDC provider, JWT callbacks |
| `packages/auth/src/middleware.ts` | Exportable auth middleware for route protection |
| `packages/auth/src/hooks.ts` | useAuth, useTenant, usePermission hooks |
| `packages/auth/src/types.ts` | Extended session/JWT types with tenant and roles |

### Modified files (each app)

For each of `apps/web-professional`, `apps/web-patient`, `apps/backoffice`:

| File | Change |
|------|--------|
| `package.json` | Add workspace deps, scripts, Tailwind deps |
| `tsconfig.json` | Extend base, add package paths |
| `next.config.mjs` | Create — transpile workspace packages |
| `tailwind.config.ts` | Create — use shared preset |
| `postcss.config.mjs` | Create — Tailwind + autoprefixer |
| `src/app/globals.css` | Create — Tailwind directives |
| `src/app/layout.tsx` | Rewrite — add providers (QueryClient, auth session, Tailwind) |
| `src/app/page.tsx` | Rewrite — basic dashboard shell with AppShell layout |
| `src/app/api/auth/[...nextauth]/route.ts` | Create — NextAuth API route |
| `src/middleware.ts` | Create — Auth middleware |

### Updated root files

| File | Change |
|------|--------|
| `.env.example` | Add frontend env vars |
| `Makefile` | Add frontend targets (dev-frontend, build-frontend, install) |

---

## Task 1: Root Monorepo Configuration

**Files:**
- Create: `package.json` (root)
- Create: `pnpm-workspace.yaml`
- Create: `turbo.json`
- Create: `tsconfig.base.json`
- Create: `.npmrc`
- Create: `.node-version`

- [ ] **Step 1: Create `.node-version`**

```
20
```

- [ ] **Step 2: Create `.npmrc`**

```ini
strict-peer-dependencies=false
auto-install-peers=true
```

- [ ] **Step 3: Create `pnpm-workspace.yaml`**

```yaml
packages:
  - "apps/*"
  - "packages/*"
```

- [ ] **Step 4: Create root `package.json`**

```json
{
  "name": "nutrometra",
  "private": true,
  "scripts": {
    "dev": "turbo dev",
    "build": "turbo build",
    "lint": "turbo lint",
    "dev:api": "cd services/api && go run ./cmd/api",
    "dev:frontend": "turbo dev --filter='./apps/*'"
  },
  "devDependencies": {
    "turbo": "^2"
  },
  "packageManager": "pnpm@9.15.4"
}
```

- [ ] **Step 5: Create `turbo.json`**

```json
{
  "$schema": "https://turbo.build/schema.json",
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": [".next/**", "dist/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "lint": {
      "dependsOn": ["^build"]
    }
  }
}
```

- [ ] **Step 6: Create `tsconfig.base.json`**

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true
  },
  "exclude": ["node_modules"]
}
```

- [ ] **Step 7: Run pnpm install at root**

```bash
pnpm install
```

Expected: Creates `pnpm-lock.yaml`, installs turbo in root `node_modules/`.

- [ ] **Step 8: Commit**

```bash
git add package.json pnpm-workspace.yaml turbo.json tsconfig.base.json .npmrc .node-version pnpm-lock.yaml
git commit -m "feat(frontend): initialize Turborepo + pnpm monorepo"
```

---

## Task 2: Package @nutrometra/ui — Foundation

**Files:**
- Create: `packages/ui/package.json`
- Create: `packages/ui/tsconfig.json`
- Create: `packages/ui/src/index.ts`
- Create: `packages/ui/src/tailwind/preset.ts`
- Create: `packages/ui/src/lib/utils.ts`

- [ ] **Step 1: Create `packages/ui/package.json`**

```json
{
  "name": "@nutrometra/ui",
  "version": "0.0.1",
  "private": true,
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "exports": {
    ".": "./src/index.ts",
    "./tailwind": "./src/tailwind/preset.ts",
    "./globals.css": "./src/globals.css"
  },
  "dependencies": {
    "class-variance-authority": "^0.7",
    "clsx": "^2",
    "tailwind-merge": "^2",
    "lucide-react": "^0.460"
  },
  "peerDependencies": {
    "react": "^18",
    "react-dom": "^18"
  },
  "devDependencies": {
    "@types/react": "^18",
    "typescript": "^5",
    "tailwindcss": "^3.4"
  }
}
```

- [ ] **Step 2: Create `packages/ui/tsconfig.json`**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "react-jsx",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.tsx"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 3: Create `packages/ui/src/tailwind/preset.ts`**

```typescript
import type { Config } from "tailwindcss"

const preset: Partial<Config> = {
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
    },
  },
}

export default preset
```

- [ ] **Step 4: Create `packages/ui/src/lib/utils.ts`**

```typescript
import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

- [ ] **Step 5: Create `packages/ui/src/globals.css`**

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --card: 0 0% 100%;
    --card-foreground: 222.2 84% 4.9%;
    --primary: 142.1 76.2% 36.3%;
    --primary-foreground: 355.7 100% 97.3%;
    --secondary: 210 40% 96%;
    --secondary-foreground: 222.2 47.4% 11.2%;
    --muted: 210 40% 96%;
    --muted-foreground: 215.4 16.3% 46.9%;
    --accent: 210 40% 96%;
    --accent-foreground: 222.2 47.4% 11.2%;
    --destructive: 0 84.2% 60.2%;
    --destructive-foreground: 210 40% 98%;
    --border: 214.3 31.8% 91.4%;
    --input: 214.3 31.8% 91.4%;
    --ring: 142.1 76.2% 36.3%;
    --radius: 0.5rem;
  }

  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    --card: 222.2 84% 4.9%;
    --card-foreground: 210 40% 98%;
    --primary: 142.1 76.2% 36.3%;
    --primary-foreground: 355.7 100% 97.3%;
    --secondary: 217.2 32.6% 17.5%;
    --secondary-foreground: 210 40% 98%;
    --muted: 217.2 32.6% 17.5%;
    --muted-foreground: 215 20.2% 65.1%;
    --accent: 217.2 32.6% 17.5%;
    --accent-foreground: 210 40% 98%;
    --destructive: 0 62.8% 30.6%;
    --destructive-foreground: 210 40% 98%;
    --border: 217.2 32.6% 17.5%;
    --input: 217.2 32.6% 17.5%;
    --ring: 142.1 76.2% 36.3%;
  }
}

@layer base {
  * {
    @apply border-border;
  }
  body {
    @apply bg-background text-foreground;
  }
}
```

- [ ] **Step 6: Create `packages/ui/src/index.ts`**

```typescript
// Utils
export { cn } from "./lib/utils"

// Layouts (will be added in Task 3)
```

- [ ] **Step 7: Run pnpm install**

```bash
pnpm install
```

- [ ] **Step 8: Commit**

```bash
git add packages/ui/
git commit -m "feat(frontend): add @nutrometra/ui package with Tailwind preset and design tokens"
```

---

## Task 3: Package @nutrometra/ui — Components and Layouts

**Files:**
- Create: `packages/ui/src/components/button.tsx`
- Create: `packages/ui/src/components/input.tsx`
- Create: `packages/ui/src/components/card.tsx`
- Create: `packages/ui/src/components/skeleton.tsx`
- Create: `packages/ui/src/layouts/app-shell.tsx`
- Create: `packages/ui/src/layouts/page-header.tsx`
- Create: `packages/ui/src/layouts/empty-state.tsx`
- Create: `packages/ui/src/layouts/loading-state.tsx`
- Create: `packages/ui/src/layouts/error-state.tsx`
- Modify: `packages/ui/src/index.ts`

- [ ] **Step 1: Create `packages/ui/src/components/button.tsx`**

Standard shadcn/ui Button component with variants (default, destructive, outline, secondary, ghost, link) and sizes (default, sm, lg, icon). Use `class-variance-authority` for variant handling and `cn()` for class merging. Include `forwardRef` and `asChild` via Radix Slot.

Dependency: add `@radix-ui/react-slot` to `packages/ui/package.json` dependencies.

- [ ] **Step 2: Create `packages/ui/src/components/input.tsx`**

Standard shadcn/ui Input component — forwardRef, accepts all HTML input props, styled with Tailwind classes.

- [ ] **Step 3: Create `packages/ui/src/components/card.tsx`**

shadcn/ui Card with sub-components: Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter.

- [ ] **Step 4: Create `packages/ui/src/components/skeleton.tsx`**

shadcn/ui Skeleton — animated pulse placeholder for loading states.

- [ ] **Step 5: Create `packages/ui/src/layouts/app-shell.tsx`**

```typescript
"use client"

import * as React from "react"
import { cn } from "../lib/utils"

interface AppShellProps {
  sidebar?: React.ReactNode
  header?: React.ReactNode
  children: React.ReactNode
  className?: string
}

export function AppShell({ sidebar, header, children, className }: AppShellProps) {
  const [sidebarOpen, setSidebarOpen] = React.useState(false)

  return (
    <div className="min-h-screen flex">
      {/* Sidebar — hidden on mobile, shown on md+ */}
      {sidebar && (
        <>
          {/* Mobile overlay */}
          {sidebarOpen && (
            <div
              className="fixed inset-0 z-40 bg-black/50 md:hidden"
              onClick={() => setSidebarOpen(false)}
            />
          )}
          <aside
            className={cn(
              "fixed inset-y-0 left-0 z-50 w-64 bg-card border-r transform transition-transform md:relative md:translate-x-0",
              sidebarOpen ? "translate-x-0" : "-translate-x-full"
            )}
          >
            {sidebar}
          </aside>
        </>
      )}

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        {header && (
          <header className="sticky top-0 z-30 border-b bg-background/95 backdrop-blur">
            <div className="flex items-center h-14 px-4 gap-4">
              {sidebar && (
                <button
                  className="md:hidden p-2 -ml-2"
                  onClick={() => setSidebarOpen(true)}
                  aria-label="Abrir menu"
                >
                  <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                  </svg>
                </button>
              )}
              {header}
            </div>
          </header>
        )}
        <main className={cn("flex-1 p-4 md:p-6", className)}>
          {children}
        </main>
      </div>
    </div>
  )
}
```

- [ ] **Step 6: Create `packages/ui/src/layouts/page-header.tsx`**

```typescript
import * as React from "react"

interface PageHeaderProps {
  title: string
  description?: string
  actions?: React.ReactNode
}

export function PageHeader({ title, description, actions }: PageHeaderProps) {
  return (
    <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
        {description && (
          <p className="text-muted-foreground text-sm mt-1">{description}</p>
        )}
      </div>
      {actions && <div className="flex gap-2 mt-2 sm:mt-0">{actions}</div>}
    </div>
  )
}
```

- [ ] **Step 7: Create layout states (empty, loading, error)**

`packages/ui/src/layouts/empty-state.tsx`:
```typescript
import * as React from "react"

interface EmptyStateProps {
  icon?: React.ReactNode
  title: string
  description?: string
  action?: React.ReactNode
}

export function EmptyState({ icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      {icon && <div className="text-muted-foreground mb-4">{icon}</div>}
      <h3 className="text-lg font-semibold">{title}</h3>
      {description && <p className="text-muted-foreground mt-1 max-w-sm">{description}</p>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}
```

`packages/ui/src/layouts/loading-state.tsx`:
```typescript
import { Skeleton } from "../components/skeleton"

interface LoadingStateProps {
  lines?: number
}

export function LoadingState({ lines = 3 }: LoadingStateProps) {
  return (
    <div className="space-y-3 py-4">
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton key={i} className="h-4 w-full" />
      ))}
    </div>
  )
}
```

`packages/ui/src/layouts/error-state.tsx`:
```typescript
interface ErrorStateProps {
  title?: string
  message: string
  onRetry?: () => void
}

export function ErrorState({ title = "Erro", message, onRetry }: ErrorStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      <div className="rounded-full bg-destructive/10 p-3 mb-4">
        <svg className="h-6 w-6 text-destructive" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
      </div>
      <h3 className="text-lg font-semibold">{title}</h3>
      <p className="text-muted-foreground mt-1">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="mt-4 text-sm font-medium text-primary hover:underline"
        >
          Tentar novamente
        </button>
      )}
    </div>
  )
}
```

- [ ] **Step 8: Update `packages/ui/src/index.ts` with all exports**

```typescript
// Utils
export { cn } from "./lib/utils"

// Components
export { Button, buttonVariants } from "./components/button"
export { Input } from "./components/input"
export { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "./components/card"
export { Skeleton } from "./components/skeleton"

// Layouts
export { AppShell } from "./layouts/app-shell"
export { PageHeader } from "./layouts/page-header"
export { EmptyState } from "./layouts/empty-state"
export { LoadingState } from "./layouts/loading-state"
export { ErrorState } from "./layouts/error-state"
```

- [ ] **Step 9: Add @radix-ui/react-slot to deps and run pnpm install**

```bash
cd packages/ui && pnpm add @radix-ui/react-slot
```

- [ ] **Step 10: Commit**

```bash
git add packages/ui/
git commit -m "feat(frontend): add shadcn/ui components and layout shells to @nutrometra/ui"
```

---

## Task 4: Package @nutrometra/api-client

**Files:**
- Create: `packages/api-client/package.json`
- Create: `packages/api-client/tsconfig.json`
- Create: `packages/api-client/src/client.ts`
- Create: `packages/api-client/src/provider.tsx`
- Create: `packages/api-client/src/types/common.ts`
- Create: `packages/api-client/src/types/auth.ts`
- Create: `packages/api-client/src/types/billing.ts`
- Create: `packages/api-client/src/types/index.ts`
- Create: `packages/api-client/src/hooks/billing.ts`
- Create: `packages/api-client/src/hooks/index.ts`
- Create: `packages/api-client/src/index.ts`

- [ ] **Step 1: Create `packages/api-client/package.json`**

```json
{
  "name": "@nutrometra/api-client",
  "version": "0.0.1",
  "private": true,
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "exports": {
    ".": "./src/index.ts",
    "./types": "./src/types/index.ts",
    "./hooks": "./src/hooks/index.ts"
  },
  "dependencies": {
    "@tanstack/react-query": "^5"
  },
  "peerDependencies": {
    "react": "^18",
    "react-dom": "^18"
  },
  "devDependencies": {
    "@types/react": "^18",
    "typescript": "^5"
  }
}
```

- [ ] **Step 2: Create `packages/api-client/tsconfig.json`**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "react-jsx"
  },
  "include": ["src/**/*.ts", "src/**/*.tsx"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 3: Create `packages/api-client/src/types/common.ts`**

```typescript
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
```

- [ ] **Step 4: Create `packages/api-client/src/types/auth.ts`**

```typescript
import type { UUID } from "./common"

export interface User {
  id: UUID
  external_id: string
  email: string
  display_name: string
  created_at: string
}

export interface TenantMembership {
  tenant_id: UUID
  tenant_name: string
  roles: string[]
}
```

- [ ] **Step 5: Create `packages/api-client/src/types/billing.ts`**

```typescript
import type { UUID } from "./common"

export interface Plan {
  id: UUID
  code: string
  name: string
  billing_cycle: string
  currency: string
  price_cents: number
}

export interface Subscription {
  id: UUID
  tenant_id: UUID
  plan_id: UUID
  status: "trialing" | "active" | "past_due" | "cancelled" | "expired"
  started_at: string
  trial_ends_at?: string
  renews_at?: string
}

export interface Entitlement {
  feature_key: string
  enabled: boolean
  limit?: number
  source: "override" | "plan" | "default"
}
```

- [ ] **Step 6: Create `packages/api-client/src/types/index.ts`**

```typescript
export type { UUID, ErrorResponse, PaginatedResponse } from "./common"
export type { User, TenantMembership } from "./auth"
export type { Plan, Subscription, Entitlement } from "./billing"
```

- [ ] **Step 7: Create `packages/api-client/src/client.ts`**

```typescript
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
}
```

- [ ] **Step 8: Create `packages/api-client/src/provider.tsx`**

```typescript
"use client"

import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import * as React from "react"

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000, // 30 seconds
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

export function ApiProvider({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )
}
```

- [ ] **Step 9: Create `packages/api-client/src/hooks/billing.ts`**

```typescript
import { useQuery } from "@tanstack/react-query"
import { api } from "../client"
import type { Plan, Subscription, Entitlement } from "../types/billing"

export function usePlans() {
  return useQuery({
    queryKey: ["plans"],
    queryFn: () => api.get<Plan[]>("/plans"),
  })
}

export function useSubscription() {
  return useQuery({
    queryKey: ["subscription"],
    queryFn: () => api.get<Subscription>("/subscription"),
  })
}

export function useEntitlements() {
  return useQuery({
    queryKey: ["entitlements"],
    queryFn: () => api.get<Record<string, Entitlement>>("/entitlements"),
    staleTime: 5 * 60 * 1000, // 5 min — aligned with backend cache TTL
  })
}
```

- [ ] **Step 10: Create `packages/api-client/src/hooks/index.ts`**

```typescript
export { usePlans, useSubscription, useEntitlements } from "./billing"
```

- [ ] **Step 11: Create `packages/api-client/src/index.ts`**

```typescript
export { api, ApiError, configureApiClient } from "./client"
export { ApiProvider } from "./provider"
export type { UUID, ErrorResponse, PaginatedResponse } from "./types/common"
export type { User, TenantMembership } from "./types/auth"
export type { Plan, Subscription, Entitlement } from "./types/billing"
```

- [ ] **Step 12: Run pnpm install and commit**

```bash
pnpm install
git add packages/api-client/
git commit -m "feat(frontend): add @nutrometra/api-client with fetch wrapper, types, and TanStack Query hooks"
```

---

## Task 5: Package @nutrometra/auth

**Files:**
- Create: `packages/auth/package.json`
- Create: `packages/auth/tsconfig.json`
- Create: `packages/auth/src/types.ts`
- Create: `packages/auth/src/config.ts`
- Create: `packages/auth/src/middleware.ts`
- Create: `packages/auth/src/hooks.ts`
- Create: `packages/auth/src/index.ts`

- [ ] **Step 1: Create `packages/auth/package.json`**

```json
{
  "name": "@nutrometra/auth",
  "version": "0.0.1",
  "private": true,
  "main": "./src/index.ts",
  "types": "./src/index.ts",
  "dependencies": {
    "next-auth": "^4"
  },
  "peerDependencies": {
    "next": "^14",
    "react": "^18"
  },
  "devDependencies": {
    "@types/react": "^18",
    "typescript": "^5"
  }
}
```

- [ ] **Step 2: Create `packages/auth/tsconfig.json`**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "react-jsx"
  },
  "include": ["src/**/*.ts", "src/**/*.tsx"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 3: Create `packages/auth/src/types.ts`**

```typescript
import type { DefaultSession, DefaultJWT } from "next-auth"

declare module "next-auth" {
  interface Session extends DefaultSession {
    accessToken?: string
    tenantId?: string
    roles?: string[]
  }
}

declare module "next-auth/jwt" {
  interface JWT extends DefaultJWT {
    accessToken?: string
    refreshToken?: string
    expiresAt?: number
    tenantId?: string
    roles?: string[]
  }
}

export type {} // ensure this is treated as a module
```

- [ ] **Step 4: Create `packages/auth/src/config.ts`**

```typescript
import type { NextAuthOptions } from "next-auth"

export function createAuthOptions(overrides?: Partial<NextAuthOptions>): NextAuthOptions {
  return {
    providers: [
      {
        id: "zitadel",
        name: "Zitadel",
        type: "oauth",
        wellKnown: `${process.env.ZITADEL_ISSUER}/.well-known/openid-configuration`,
        clientId: process.env.ZITADEL_CLIENT_ID,
        clientSecret: process.env.ZITADEL_CLIENT_SECRET,
        authorization: {
          params: {
            scope: "openid profile email",
          },
        },
        idToken: true,
        profile(profile) {
          return {
            id: profile.sub,
            name: profile.name ?? profile.preferred_username,
            email: profile.email,
          }
        },
      },
    ],
    callbacks: {
      async jwt({ token, account }) {
        if (account) {
          token.accessToken = account.access_token
          token.refreshToken = account.refresh_token
          token.expiresAt = account.expires_at
        }
        return token
      },
      async session({ session, token }) {
        session.accessToken = token.accessToken as string | undefined
        session.tenantId = token.tenantId as string | undefined
        session.roles = token.roles as string[] | undefined
        return session
      },
    },
    pages: {
      signIn: "/auth/signin",
    },
    session: {
      strategy: "jwt",
    },
    ...overrides,
  }
}
```

- [ ] **Step 5: Create `packages/auth/src/middleware.ts`**

```typescript
import { withAuth } from "next-auth/middleware"

export function createAuthMiddleware(publicPaths: string[] = []) {
  return withAuth({
    pages: {
      signIn: "/auth/signin",
    },
    callbacks: {
      authorized({ token }) {
        return !!token
      },
    },
  })
}

export function createMiddlewareMatcher(publicPaths: string[] = []) {
  return {
    matcher: [
      "/((?!api/auth|_next/static|_next/image|favicon.ico|auth).*)",
      ...publicPaths,
    ],
  }
}
```

- [ ] **Step 6: Create `packages/auth/src/hooks.ts`**

```typescript
"use client"

import { useSession } from "next-auth/react"

export function useAuth() {
  const { data: session, status } = useSession()
  return {
    user: session?.user ?? null,
    accessToken: session?.accessToken ?? null,
    isAuthenticated: status === "authenticated",
    isLoading: status === "loading",
    tenantId: session?.tenantId ?? null,
    roles: session?.roles ?? [],
  }
}

export function useTenant() {
  const { tenantId } = useAuth()
  return { tenantId }
}

export function usePermission(code: string): boolean {
  // For now, return true — will be wired to useEntitlements in portal implementation
  return true
}
```

- [ ] **Step 7: Create `packages/auth/src/index.ts`**

```typescript
export { createAuthOptions } from "./config"
export { createAuthMiddleware, createMiddlewareMatcher } from "./middleware"
export { useAuth, useTenant, usePermission } from "./hooks"
export type {} from "./types" // side-effect: module augmentation
```

- [ ] **Step 8: Run pnpm install and commit**

```bash
pnpm install
git add packages/auth/
git commit -m "feat(frontend): add @nutrometra/auth with next-auth Zitadel OIDC config and hooks"
```

---

## Task 6: Wire web-professional App

This task wires one app completely. Tasks 7 and 8 replicate for the other two apps.

**Files:**
- Modify: `apps/web-professional/package.json`
- Modify: `apps/web-professional/tsconfig.json`
- Create: `apps/web-professional/next.config.mjs`
- Create: `apps/web-professional/tailwind.config.ts`
- Create: `apps/web-professional/postcss.config.mjs`
- Create: `apps/web-professional/src/app/globals.css`
- Modify: `apps/web-professional/src/app/layout.tsx`
- Modify: `apps/web-professional/src/app/page.tsx`
- Create: `apps/web-professional/src/app/api/auth/[...nextauth]/route.ts`
- Create: `apps/web-professional/src/middleware.ts`

- [ ] **Step 1: Update `apps/web-professional/package.json`**

```json
{
  "name": "web-professional",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev -p 3000",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "@nutrometra/api-client": "workspace:*",
    "@nutrometra/auth": "workspace:*",
    "@nutrometra/ui": "workspace:*",
    "next": "14.2.5",
    "next-auth": "^4",
    "react": "^18",
    "react-dom": "^18",
    "@tanstack/react-query": "^5"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "autoprefixer": "^10",
    "postcss": "^8",
    "tailwindcss": "^3.4",
    "typescript": "^5"
  }
}
```

- [ ] **Step 2: Update `apps/web-professional/tsconfig.json`**

```json
{
  "extends": "../../tsconfig.base.json",
  "compilerOptions": {
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./src/*"],
      "@nutrometra/ui": ["../../packages/ui/src"],
      "@nutrometra/api-client": ["../../packages/api-client/src"],
      "@nutrometra/auth": ["../../packages/auth/src"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

- [ ] **Step 3: Create `apps/web-professional/next.config.mjs`**

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  transpilePackages: [
    "@nutrometra/ui",
    "@nutrometra/api-client",
    "@nutrometra/auth",
  ],
}

export default nextConfig
```

- [ ] **Step 4: Create `apps/web-professional/tailwind.config.ts`**

```typescript
import type { Config } from "tailwindcss"
import sharedPreset from "@nutrometra/ui/tailwind"

const config: Config = {
  presets: [sharedPreset as Config],
  content: [
    "./src/**/*.{ts,tsx}",
    "../../packages/ui/src/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

export default config
```

- [ ] **Step 5: Create `apps/web-professional/postcss.config.mjs`**

```javascript
const config = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}

export default config
```

- [ ] **Step 6: Create `apps/web-professional/src/app/globals.css`**

```css
@import "@nutrometra/ui/globals.css";
```

If direct import doesn't work with the package setup, copy the CSS content from `packages/ui/src/globals.css` instead. The Tailwind directives and CSS variables must be present.

- [ ] **Step 7: Create `apps/web-professional/src/app/api/auth/[...nextauth]/route.ts`**

```typescript
import NextAuth from "next-auth"
import { createAuthOptions } from "@nutrometra/auth"

const handler = NextAuth(createAuthOptions())

export { handler as GET, handler as POST }
```

- [ ] **Step 8: Create `apps/web-professional/src/middleware.ts`**

```typescript
export { default } from "next-auth/middleware"

export const config = {
  matcher: ["/((?!api/auth|_next/static|_next/image|favicon.ico|auth).*)"],
}
```

- [ ] **Step 9: Rewrite `apps/web-professional/src/app/layout.tsx`**

```typescript
import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { ApiProvider } from "@nutrometra/api-client"
import { SessionProvider } from "next-auth/react"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "Nutrometra — Portal Profissional",
  description: "Plataforma de nutrição para profissionais",
}

function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <ApiProvider>
        {children}
      </ApiProvider>
    </SessionProvider>
  )
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="pt-BR">
      <body className={inter.className}>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  )
}
```

Note: `SessionProvider` must be in a `"use client"` wrapper. Extract `Providers` into a separate `src/app/providers.tsx` file with `"use client"` directive if needed.

- [ ] **Step 10: Rewrite `apps/web-professional/src/app/page.tsx`**

```typescript
import { AppShell, PageHeader } from "@nutrometra/ui"

function Sidebar() {
  return (
    <nav className="flex flex-col gap-1 p-4">
      <h2 className="text-lg font-semibold mb-4 px-2">Nutrometra</h2>
      <span className="text-sm text-muted-foreground px-2">Menu em breve</span>
    </nav>
  )
}

export default function HomePage() {
  return (
    <AppShell sidebar={<Sidebar />} header={<span className="font-semibold">Portal Profissional</span>}>
      <PageHeader title="Dashboard" description="Bem-vindo ao Nutrometra" />
      <p className="text-muted-foreground">Portal profissional em construcao.</p>
    </AppShell>
  )
}
```

- [ ] **Step 11: Run pnpm install and verify build**

```bash
pnpm install
pnpm turbo build --filter=web-professional
```

Expected: Build completes successfully.

- [ ] **Step 12: Commit**

```bash
git add apps/web-professional/
git commit -m "feat(frontend): wire web-professional with shared packages, Tailwind, and auth"
```

---

## Task 7: Wire web-patient App

Same structure as Task 6 but for `apps/web-patient/`. Key differences:
- Port: 3001
- Title: "Nutrometra — Portal Paciente"
- Page header: "Portal Paciente"
- Simpler layout (no sidebar initially — mobile-first)

**Files:** Same pattern as Task 6, all under `apps/web-patient/`.

- [ ] **Steps 1-12:** Replicate Task 6 for `apps/web-patient/` with these changes:
  - `package.json` name: `"web-patient"`, port: `-p 3001`
  - `layout.tsx` metadata title: `"Nutrometra — Portal Paciente"`
  - `page.tsx`: No sidebar, just a centered card layout for mobile-first:

```typescript
import { PageHeader, Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"

export default function HomePage() {
  return (
    <div className="min-h-screen p-4 max-w-lg mx-auto">
      <PageHeader title="Portal Paciente" />
      <Card>
        <CardHeader>
          <CardTitle>Bem-vindo</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">Portal do paciente em construcao.</p>
        </CardContent>
      </Card>
    </div>
  )
}
```

- [ ] **Commit:**

```bash
git add apps/web-patient/
git commit -m "feat(frontend): wire web-patient with shared packages, Tailwind, and auth"
```

---

## Task 8: Wire backoffice App

Same structure as Task 6 but for `apps/backoffice/`. Key differences:
- Port: 3002
- Title: "Nutrometra — Backoffice"
- Page header: "Backoffice"
- Uses sidebar layout (table-heavy admin UI)

**Files:** Same pattern as Task 6, all under `apps/backoffice/`.

- [ ] **Steps 1-12:** Replicate Task 6 for `apps/backoffice/` with:
  - `package.json` name: `"backoffice"`, port: `-p 3002`
  - `layout.tsx` metadata title: `"Nutrometra — Backoffice"`
  - `page.tsx` with sidebar:

```typescript
import { AppShell, PageHeader } from "@nutrometra/ui"

function Sidebar() {
  return (
    <nav className="flex flex-col gap-1 p-4">
      <h2 className="text-lg font-semibold mb-4 px-2">Backoffice</h2>
      <span className="text-sm text-muted-foreground px-2">Menu em breve</span>
    </nav>
  )
}

export default function HomePage() {
  return (
    <AppShell sidebar={<Sidebar />} header={<span className="font-semibold">Backoffice</span>}>
      <PageHeader title="Dashboard" description="Painel administrativo" />
      <p className="text-muted-foreground">Backoffice em construcao.</p>
    </AppShell>
  )
}
```

- [ ] **Commit:**

```bash
git add apps/backoffice/
git commit -m "feat(frontend): wire backoffice with shared packages, Tailwind, and auth"
```

---

## Task 9: Root Config Updates and Verification

**Files:**
- Modify: `.env.example`
- Modify: `Makefile`

- [ ] **Step 1: Update `.env.example`** — add frontend section:

```bash
# Frontend — API
NEXT_PUBLIC_API_URL=http://localhost:8081

# Frontend — Auth
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=
ZITADEL_CLIENT_ID=
ZITADEL_CLIENT_SECRET=
ZITADEL_ISSUER=
```

Do NOT put actual values. Fields for secrets are left empty.

- [ ] **Step 2: Add Makefile frontend targets:**

```makefile
# Frontend
install-frontend:
	pnpm install

dev-frontend:
	pnpm turbo dev --filter='./apps/*'

build-frontend:
	pnpm turbo build --filter='./apps/*'
```

- [ ] **Step 3: Full verification**

```bash
pnpm install
pnpm turbo build
```

Expected: All 3 packages compile, all 3 apps build successfully.

- [ ] **Step 4: Start dev to verify**

```bash
pnpm turbo dev --filter=web-professional
```

Expected: App starts on http://localhost:3000, redirects to auth/signin (no Zitadel running = expected 404 on OIDC discovery, confirming auth middleware works).

- [ ] **Step 5: Commit**

```bash
git add .env.example Makefile pnpm-lock.yaml
git commit -m "feat(frontend): add env vars documentation and Makefile frontend targets"
```

---

## Verification Checklist

After all tasks are complete:

1. `pnpm install` — resolves all workspace dependencies
2. `pnpm turbo build` — all 3 packages + 3 apps compile
3. Each app shows styled page with Tailwind CSS applied (green primary color from theme)
4. AppShell responsive — sidebar collapses on mobile viewport
5. Auth middleware redirects unauthenticated users to `/auth/signin`
6. TypeScript compiles without errors across all packages
7. TanStack Query provider is initialized in each app's layout
