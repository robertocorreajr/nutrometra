import { describe, it, expect, afterEach } from "vitest"
import { render, screen, cleanup } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { DataTable, type Column } from "../components/data-table"

afterEach(cleanup)

interface TestItem { id: string; name: string; email: string }

const testData: TestItem[] = [
  { id: "1", name: "Alice", email: "alice@test.com" },
  { id: "2", name: "Bob", email: "bob@test.com" },
  { id: "3", name: "Charlie", email: "charlie@test.com" },
]

const columns: Column<TestItem>[] = [
  { key: "name", header: "Nome", render: (i) => i.name, searchValue: (i) => i.name, sortable: true },
  { key: "email", header: "E-mail", render: (i) => i.email, searchValue: (i) => i.email },
]

describe("DataTable", () => {
  it("renders all items", () => {
    render(<DataTable data={testData} columns={columns} keyExtractor={(i) => i.id} />)
    expect(screen.getAllByText("Alice").length).toBeGreaterThan(0)
    expect(screen.getAllByText("Bob").length).toBeGreaterThan(0)
  })

  it("filters by search term", async () => {
    const user = userEvent.setup()
    render(<DataTable data={testData} columns={columns} keyExtractor={(i) => i.id} />)
    // There is one search input rendered (toolbar area); getAllByPlaceholderText is safe here
    const [searchInput] = screen.getAllByPlaceholderText("Buscar...")
    await user.type(searchInput, "alice")
    expect(screen.getAllByText("Alice").length).toBeGreaterThan(0)
    expect(screen.queryAllByText("Bob")).toHaveLength(0)
  })

  it("shows empty message", () => {
    render(<DataTable data={[]} columns={columns} keyExtractor={(i) => i.id} emptyMessage="Vazio." />)
    expect(screen.getAllByText("Vazio.").length).toBeGreaterThan(0)
  })
})
