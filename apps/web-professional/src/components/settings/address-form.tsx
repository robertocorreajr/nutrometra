"use client"

import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Button, Input, Card, CardContent, CardHeader, CardTitle } from "@nutrometra/ui"
import {
  addressSchema,
  type AddressFormValues,
  brazilianStates,
} from "@/lib/schemas/professional"

interface AddressFormProps {
  onSubmit: (data: AddressFormValues) => void
  isSubmitting: boolean
  onCancel: () => void
  defaultValues?: Partial<AddressFormValues>
}

export function AddressForm({ onSubmit, isSubmitting, onCancel, defaultValues }: AddressFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<AddressFormValues>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      label: defaultValues?.label ?? "",
      street: defaultValues?.street ?? "",
      number: defaultValues?.number ?? "",
      complement: defaultValues?.complement ?? "",
      neighborhood: defaultValues?.neighborhood ?? "",
      city: defaultValues?.city ?? "",
      state: defaultValues?.state ?? "",
      zip_code: defaultValues?.zip_code ?? "",
      country: defaultValues?.country ?? "BR",
      phone: defaultValues?.phone ?? "",
      notes: defaultValues?.notes ?? "",
    },
  })

  const selectClasses =
    "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"

  return (
    <Card>
      <CardHeader>
        <CardTitle>{defaultValues ? "Editar Endereço" : "Novo Endereço"}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-1">
            <label htmlFor="label" className="text-sm font-medium">
              Nome do endereço
            </label>
            <Input
              id="label"
              {...register("label")}
              placeholder="Ex: Consultório Centro, Clínica Zona Sul"
            />
            {errors.label && (
              <p className="text-sm text-destructive">{errors.label.message}</p>
            )}
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="sm:col-span-2 space-y-1">
              <label htmlFor="street" className="text-sm font-medium">
                Rua
              </label>
              <Input id="street" {...register("street")} placeholder="Rua, Avenida, etc." />
              {errors.street && (
                <p className="text-sm text-destructive">{errors.street.message}</p>
              )}
            </div>

            <div className="space-y-1">
              <label htmlFor="number" className="text-sm font-medium">
                Número
              </label>
              <Input id="number" {...register("number")} placeholder="123" />
              {errors.number && (
                <p className="text-sm text-destructive">{errors.number.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="space-y-1">
              <label htmlFor="complement" className="text-sm font-medium">
                Complemento
              </label>
              <Input id="complement" {...register("complement")} placeholder="Sala 101, Bloco A" />
              {errors.complement && (
                <p className="text-sm text-destructive">{errors.complement.message}</p>
              )}
            </div>

            <div className="space-y-1">
              <label htmlFor="neighborhood" className="text-sm font-medium">
                Bairro
              </label>
              <Input id="neighborhood" {...register("neighborhood")} placeholder="Bairro" />
              {errors.neighborhood && (
                <p className="text-sm text-destructive">{errors.neighborhood.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div className="space-y-1">
              <label htmlFor="city" className="text-sm font-medium">
                Cidade
              </label>
              <Input id="city" {...register("city")} placeholder="Cidade" />
              {errors.city && (
                <p className="text-sm text-destructive">{errors.city.message}</p>
              )}
            </div>

            <div className="space-y-1">
              <label htmlFor="state" className="text-sm font-medium">
                Estado
              </label>
              <select id="state" {...register("state")} className={selectClasses}>
                <option value="">Selecione...</option>
                {brazilianStates.map((uf) => (
                  <option key={uf} value={uf}>
                    {uf}
                  </option>
                ))}
              </select>
              {errors.state && (
                <p className="text-sm text-destructive">{errors.state.message}</p>
              )}
            </div>

            <div className="space-y-1">
              <label htmlFor="zip_code" className="text-sm font-medium">
                CEP
              </label>
              <Input id="zip_code" {...register("zip_code")} placeholder="00000-000" />
              {errors.zip_code && (
                <p className="text-sm text-destructive">{errors.zip_code.message}</p>
              )}
            </div>
          </div>

          <div className="space-y-1">
            <label htmlFor="phone" className="text-sm font-medium">
              Telefone do local
            </label>
            <Input id="phone" {...register("phone")} placeholder="(11) 3333-3333" />
            {errors.phone && (
              <p className="text-sm text-destructive">{errors.phone.message}</p>
            )}
          </div>

          <div className="space-y-1">
            <label htmlFor="notes" className="text-sm font-medium">
              Observações
            </label>
            <textarea
              id="notes"
              {...register("notes")}
              rows={3}
              placeholder="Instruções de acesso, ponto de referência, etc."
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            />
            {errors.notes && (
              <p className="text-sm text-destructive">{errors.notes.message}</p>
            )}
          </div>

          <div className="flex gap-2 justify-end">
            <Button type="button" variant="ghost" onClick={onCancel}>
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Salvando..." : "Salvar Endereço"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  )
}
