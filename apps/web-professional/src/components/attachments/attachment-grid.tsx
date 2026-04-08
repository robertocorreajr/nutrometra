import { Card, CardContent } from "@nutrometra/ui"
import type { ClinicalAttachment } from "@nutrometra/api-client"
import { format } from "date-fns"
import { ptBR } from "date-fns/locale"
import { FileText, Image, FlaskConical, Pill, File } from "lucide-react"
import { attachmentCategories } from "@/lib/schemas/attachment"

const categoryIcons: Record<string, React.ElementType> = {
  general: FileText,
  exam: FlaskConical,
  lab_result: FlaskConical,
  prescription: Pill,
  photo: Image,
  other: File,
}

interface AttachmentGridProps {
  attachments: ClinicalAttachment[]
}

export function AttachmentGrid({ attachments }: AttachmentGridProps) {
  const sorted = [...attachments].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )

  function getCategoryLabel(value: string): string {
    return attachmentCategories.find((c) => c.value === value)?.label ?? value
  }

  function formatFileSize(bytes: number): string {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {sorted.map((attachment) => {
        const Icon = categoryIcons[attachment.category] ?? File
        return (
          <Card key={attachment.id} className="hover:bg-accent/50 transition-colors">
            <CardContent className="p-4">
              <div className="flex items-start gap-3">
                <div className="rounded-lg bg-muted p-2">
                  <Icon className="h-5 w-5 text-muted-foreground" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{attachment.file_name}</p>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="text-xs text-muted-foreground">{getCategoryLabel(attachment.category)}</span>
                    <span className="text-xs text-muted-foreground">{formatFileSize(attachment.file_size_bytes)}</span>
                  </div>
                  {attachment.description && (
                    <p className="text-xs text-muted-foreground mt-1 line-clamp-2">{attachment.description}</p>
                  )}
                  <p className="text-xs text-muted-foreground mt-2">
                    {format(new Date(attachment.created_at), "dd/MM/yyyy", { locale: ptBR })}
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}
