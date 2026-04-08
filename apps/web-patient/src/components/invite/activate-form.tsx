"use client"

import { useState } from "react"
import { Button, Input } from "@nutrometra/ui"

interface ActivateFormProps {
  onActivate: (code: string) => void
  isPending: boolean
  error?: string
}

export function ActivateForm({ onActivate, isPending, error }: ActivateFormProps) {
  const [code, setCode] = useState("")

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (code.trim().length > 0) {
      onActivate(code.trim().toUpperCase())
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="invite-code" className="block text-sm font-medium mb-2">
          Codigo de convite
        </label>
        <Input
          id="invite-code"
          value={code}
          onChange={(e) => setCode(e.target.value.toUpperCase())}
          placeholder="Ex: ABC123"
          className="text-center text-2xl tracking-widest font-mono h-14"
          maxLength={10}
          autoFocus
          autoComplete="off"
        />
      </div>
      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}
      <Button
        type="submit"
        className="w-full h-12 text-base"
        disabled={isPending || code.trim().length === 0}
      >
        {isPending ? "Ativando..." : "Ativar Convite"}
      </Button>
    </form>
  )
}
