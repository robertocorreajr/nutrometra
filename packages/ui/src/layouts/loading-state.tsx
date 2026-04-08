import { Skeleton } from "../components/skeleton"

interface LoadingStateProps {
  lines?: number
}

export function LoadingState({ lines = 3 }: LoadingStateProps) {
  return (
    <div className="space-y-3 py-4">
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton key={i} className="h-4 w-full" />
      ))}
    </div>
  )
}
