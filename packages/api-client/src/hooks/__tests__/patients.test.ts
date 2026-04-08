import { describe, it, expect, vi, beforeEach } from "vitest"
import { renderHook, waitFor } from "@testing-library/react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { createElement } from "react"
import { usePatients } from "../patients"
import { api } from "../../client"

vi.mock("../../client", () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  return ({ children }: { children: React.ReactNode }) =>
    createElement(QueryClientProvider, { client: queryClient }, children)
}

describe("usePatients", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it("fetches patients list", async () => {
    const mockPatients = [
      {
        id: "p1",
        tenant_id: "t1",
        professional_id: "prof1",
        full_name: "João Silva",
        active: true,
        created_at: "2024-01-01T00:00:00Z",
        updated_at: "2024-01-01T00:00:00Z",
      },
    ]
    vi.mocked(api.get).mockResolvedValueOnce(mockPatients)

    const { result } = renderHook(() => usePatients(), {
      wrapper: createWrapper(),
    })

    await waitFor(() => expect(result.current.isSuccess).toBe(true))
    expect(result.current.data).toEqual(mockPatients)
    expect(api.get).toHaveBeenCalledWith("/patients")
  })
})
