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
