import type { UserRole } from "../types/user";

const currencyFormatter = new Intl.NumberFormat("ru-RU", {
  style: "currency",
  currency: "RUB",
  maximumFractionDigits: 2,
});

export const formatPrice = (price: number): string => currencyFormatter.format(price);

export const formatDateTime = (value: string): string => {
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return value;
  }

  return parsed.toLocaleString("ru-RU");
};

export const shortId = (id: string): string => {
  if (!id) {
    return "—";
  }

  if (id.length <= 12) {
    return id;
  }

  return `${id.slice(0, 8)}...${id.slice(-4)}`;
};

export const roleLabel = (role: UserRole): string => {
  if (role === "buyer") {
    return "Покупатель";
  }

  if (role === "seller") {
    return "Продавец";
  }

  if (role === "active") {
    return "Активный";
  }

  if (role === "blocked") {
    return "Заблокирован";
  }

  return role || "Неизвестно";
};

export const requestStatusLabel = (status: string): string => {
  if (status === "pending") {
    return "в ожидании";
  }

  if (status === "approved") {
    return "одобрено";
  }

  if (status === "rejected") {
    return "отклонено";
  }

  return status;
};
