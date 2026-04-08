"use client"

import { useAuth } from "@nutrometra/auth"
import { PageHeader, Card, CardContent, CardHeader, CardTitle, Button, LoadingState, ErrorState } from "@nutrometra/ui"
import { useMyProfile } from "@nutrometra/api-client/hooks"
import { LogOut } from "lucide-react"
import { signOut } from "next-auth/react"

export default function PerfilPage() {
  const { user, patientId } = useAuth()
  const { data: profile, isLoading, isError, refetch } = useMyProfile(patientId ?? "")

  return (
    <div className="max-w-lg mx-auto md:max-w-none space-y-4">
      <PageHeader title="Meu Perfil" />

      {/* Basic info from auth */}
      <Card>
        <CardHeader><CardTitle className="text-base">Dados Pessoais</CardTitle></CardHeader>
        <CardContent className="space-y-2">
          {user?.name && (
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Nome</span>
              <span>{user.name}</span>
            </div>
          )}
          {user?.email && (
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Email</span>
              <span>{user.email}</span>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Health profile */}
      {isLoading && <LoadingState />}
      {isError && <ErrorState message="Nao foi possivel carregar o perfil de saude." onRetry={refetch} />}
      {profile && (
        <Card>
          <CardHeader><CardTitle className="text-base">Perfil de Saude</CardTitle></CardHeader>
          <CardContent className="space-y-2">
            {profile.blood_type && (
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Tipo Sanguineo</span>
                <span>{profile.blood_type}</span>
              </div>
            )}
            {profile.allergies.length > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Alergias</span>
                <span>{profile.allergies.join(", ")}</span>
              </div>
            )}
            {profile.chronic_conditions.length > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Condicoes Cronicas</span>
                <span>{profile.chronic_conditions.join(", ")}</span>
              </div>
            )}
            {profile.medications.length > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Medicamentos</span>
                <span>{profile.medications.join(", ")}</span>
              </div>
            )}
            {profile.emergency_contact_name && (
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Contato Emergencia</span>
                <span>{profile.emergency_contact_name} {profile.emergency_contact_phone ? `(${profile.emergency_contact_phone})` : ""}</span>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Logout button — prominent on mobile for easy access */}
      <div className="pt-4">
        <Button
          variant="outline"
          className="w-full"
          onClick={() => signOut({ callbackUrl: "/signin" })}
        >
          <LogOut className="h-4 w-4 mr-2" />
          Sair da conta
        </Button>
      </div>
    </div>
  )
}
