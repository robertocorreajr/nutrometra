import { Card, CardContent } from "@nutrometra/ui"

interface WelcomeCardProps {
  name?: string
}

export function WelcomeCard({ name }: WelcomeCardProps) {
  const greeting = getGreeting()
  return (
    <Card>
      <CardContent className="pt-6">
        <h2 className="text-xl font-semibold">
          {greeting}{name ? `, ${name}` : ""}!
        </h2>
        <p className="text-sm text-muted-foreground mt-1">
          Acompanhe aqui suas dietas, consultas e medidas
        </p>
      </CardContent>
    </Card>
  )
}

function getGreeting(): string {
  const hour = new Date().getHours()
  if (hour < 12) return "Bom dia"
  if (hour < 18) return "Boa tarde"
  return "Boa noite"
}
