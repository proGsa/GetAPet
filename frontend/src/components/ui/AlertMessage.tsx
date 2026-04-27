import type { ReactNode } from "react";

interface AlertMessageProps {
  variant?: "error" | "success" | "info";
  children: ReactNode;
}

export function AlertMessage({ variant = "info", children }: AlertMessageProps) {
  return <div className={`alert alert-${variant}`}>{children}</div>;
}
