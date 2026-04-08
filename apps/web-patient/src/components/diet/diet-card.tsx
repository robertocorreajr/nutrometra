import Link from "next/link"
import { Card, CardContent } from "@nutrometra/ui"
import { format, parseISO } from "date-fns"
import type { Diet } from "@nutrometra/api-client/types"

interface DietCardProps {
  diet: Diet
}

const statusLabels: Record<string, string> = {
  draft: "Rascunho", published: "Ativa", archived: "Arquivada",
}
const statusColors: Record<string, string> = {
  draft: "bg-gray-100 text-gray-700", published: "bg-green-100 text-green-700", archived: "bg-yellow-100 text-yellow-700",
}

export function DietCard({ diet }: DietCardProps) {
  return (
    <Link href={`/dietas/${diet.id}`}>
      <Card className="hover:border-primary/50 transition-colors">
        <CardContent className="pt-4">
          <div className="flex items-start justify-between">
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold truncate">{diet.title}</h3>
              {diet.objective && <p className="text-sm text-muted-foreground mt-1 line-clamp-2">{diet.objective}</p>}
            </div>
            <span className={`ml-2 inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${statusColors[diet.status] ?? "bg-gray-100 text-gray-700"}`}>
              {statusLabels[diet.status] ?? diet.status}
            </span>
          </div>
          <div className="flex gap-4 mt-3 text-xs text-muted-foreground">
            {diet.published_at && <span>Publicada: {format(parseISO(diet.published_at), "dd/MM/yyyy")}</span>}
            {diet.valid_from && <span>Valida: {format(parseISO(diet.valid_from), "dd/MM/yy")} {diet.valid_until ? `- ${format(parseISO(diet.valid_until), "dd/MM/yy")}` : ""}</span>}
          </div>
        </CardContent>
      </Card>
    </Link>
  )
}
