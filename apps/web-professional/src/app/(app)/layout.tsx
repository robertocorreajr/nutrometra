import { AppShell } from "@nutrometra/ui"
import { Sidebar } from "@/components/sidebar"
import { Header } from "@/components/header"

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AppShell sidebar={<Sidebar />} header={<Header />}>
      {children}
    </AppShell>
  )
}
