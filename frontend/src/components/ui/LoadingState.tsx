interface LoadingStateProps {
  label?: string;
}

export function LoadingState({ label = "Загрузка..." }: LoadingStateProps) {
  return (
    <div className="loading-state" role="status" aria-live="polite">
      <span className="spinner" />
      <span>{label}</span>
    </div>
  );
}
