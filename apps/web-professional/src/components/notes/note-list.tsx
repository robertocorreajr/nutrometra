import { Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import type { ProgressNote } from "@nutrometra/api-client"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { Eye, EyeOff } from "lucide-react"

interface NoteListProps {
  notes: ProgressNote[]
}

export function NoteList({ notes }: NoteListProps) {
  const sorted = [...notes].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )

  return (
    <div className="space-y-3">
      {sorted.map((note) => (
        <Card key={note.id}>
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-sm font-medium">{note.title}</CardTitle>
              <div className="flex items-center gap-2">
                {note.visible_to_patient ? (
                  <span className="flex items-center gap-1 text-xs text-muted-foreground"><Eye className="h-3 w-3" />Visível</span>
                ) : (
                  <span className="flex items-center gap-1 text-xs text-muted-foreground"><EyeOff className="h-3 w-3" />Oculto</span>
                )}
                <span className="text-xs text-muted-foreground">
                  {format(new Date(note.created_at), "dd/MM/yyyy HH:mm", { locale: ptBR })}
                </span>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm whitespace-pre-wrap">{note.content}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
