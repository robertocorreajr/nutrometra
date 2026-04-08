"use client"

import { useRef } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent } from "@nutrometra/ui"
import { attachmentSchema, type AttachmentFormValues, attachmentCategories } from "@/lib/schemas/attachment"

interface AttachmentFormProps {
  onSubmit: (data: AttachmentFormValues) => void
  isSubmitting: boolean
  onCancel?: () => void
}

export function AttachmentForm({ onSubmit, isSubmitting, onCancel }: AttachmentFormProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<AttachmentFormValues>({
    resolver: zodResolver(attachmentSchema),
    defaultValues: {
      file_name: "",
      file_type: "",
      file_size_bytes: 0,
      category: "general",
      description: "",
    },
  })

  const fileName = watch("file_name")

  function handleFileSelect(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    setValue("file_name", file.name)
    setValue("file_type", file.type || "application/octet-stream")
    setValue("file_size_bytes", file.size)
  }

  function handleFormSubmit(data: AttachmentFormValues) {
    onSubmit(data)
    reset()
    if (fileInputRef.current) fileInputRef.current.value = ""
  }

  return (
    <Card>
      <CardContent className="pt-6">
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Arquivo</label>
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileSelect}
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm file:border-0 file:bg-transparent file:text-sm file:font-medium"
            />
            {fileName && <p className="text-xs text-muted-foreground">Selecionado: {fileName}</p>}
            {errors.file_name && <p className="text-sm text-destructive">{errors.file_name.message}</p>}
            <p className="text-xs text-muted-foreground">Nota: nesta versão, apenas os metadados do arquivo são registrados. O upload para storage será integrado em versão futura.</p>
          </div>

          <div className="space-y-2">
            <label htmlFor="category" className="text-sm font-medium">Categoria</label>
            <select
              id="category"
              {...register("category")}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              {attachmentCategories.map((cat) => (
                <option key={cat.value} value={cat.value}>{cat.label}</option>
              ))}
            </select>
          </div>

          <div className="space-y-2">
            <label htmlFor="description" className="text-sm font-medium">Descrição</label>
            <textarea
              id="description"
              {...register("description")}
              rows={3}
              placeholder="Descreva o conteúdo do anexo..."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            />
          </div>

          <div className="flex gap-2 justify-end">
            {onCancel && <Button type="button" variant="ghost" onClick={onCancel}>Cancelar</Button>}
            <Button type="submit" disabled={isSubmitting || !fileName}>{isSubmitting ? "Registrando..." : "Registrar Anexo"}</Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
