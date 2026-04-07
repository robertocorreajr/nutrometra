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
