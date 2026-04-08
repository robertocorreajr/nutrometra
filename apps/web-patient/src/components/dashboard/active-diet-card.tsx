import Link from "next/link"
import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import { UtensilsCrossed } from "lucide-react"
import type { Diet } from "@nutrometra/api-client/types"

interface ActiveDietCardProps {
  diet?: Diet
}

export function ActiveDietCard({ diet }: ActiveDietCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground flex items-center gap-2">
          <UtensilsCrossed className="h-4 w-4" />
          Dieta Ativa
        </CardTitle>
      </CardHeader>
      <CardContent>
        {diet ? (
          <Link href={`/dietas/${diet.id}`} className="block group">
            <p className="text-lg font-semibold group-hover:text-primary transition-colors">
              {diet.title}
            </p>
            {diet.objective && (
              <p className="text-xs text-muted-foreground mt-1 line-clamp-2">{diet.objective}</p>
            )}
          </Link>
        ) : (
          <p className="text-sm text-muted-foreground">Nenhuma dieta ativa</p>
        )}
      </CardContent>
    </Card>
  )
}
