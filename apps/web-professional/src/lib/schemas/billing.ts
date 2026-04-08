export const subscriptionStatusLabels: Record<string, string> = {
  trialing: "Em Trial",
  active: "Ativa",
  past_due: "Pagamento Pendente",
  cancelled: "Cancelada",
  expired: "Expirada",
}

export const subscriptionStatusColors: Record<string, string> = {
  trialing: "bg-blue-100 text-blue-800",
  active: "bg-green-100 text-green-800",
  past_due: "bg-amber-100 text-amber-800",
  cancelled: "bg-red-100 text-red-800",
  expired: "bg-gray-100 text-gray-800",
}

export const invoiceStatusLabels: Record<string, string> = {
  draft: "Rascunho",
  open: "Em Aberto",
  paid: "Paga",
  void: "Anulada",
  uncollectible: "Inadimplente",
}

export const invoiceStatusColors: Record<string, string> = {
  draft: "bg-gray-100 text-gray-800",
  open: "bg-blue-100 text-blue-800",
  paid: "bg-green-100 text-green-800",
  void: "bg-gray-100 text-gray-800",
  uncollectible: "bg-red-100 text-red-800",
}

export const paymentStatusLabels: Record<string, string> = {
  pending: "Pendente",
  succeeded: "Aprovado",
  failed: "Falhou",
  refunded: "Reembolsado",
}

export const paymentStatusColors: Record<string, string> = {
  pending: "bg-amber-100 text-amber-800",
  succeeded: "bg-green-100 text-green-800",
  failed: "bg-red-100 text-red-800",
  refunded: "bg-blue-100 text-blue-800",
}

export function formatCurrency(amountCents: number, currency: string): string {
  const amount = amountCents / 100
  if (currency === "BRL") {
    return new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(amount)
  }
  return new Intl.NumberFormat("en-US", { style: "currency", currency }).format(amount)
}
