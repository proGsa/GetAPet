import type { BackendErrorPayload } from "../types/common";

const DEFAULT_API_BASE = "/api";
const API_BASE = import.meta.env.VITE_API_BASE ?? DEFAULT_API_BASE;

export interface RequestOptions extends RequestInit {
  token?: string;
}

export class ApiError extends Error {
  public readonly status: number;
  public readonly backendError?: string;

  constructor(status: number, message: string, backendError?: string) {
    super(message);
    this.status = status;
    this.backendError = backendError;
  }
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null;

const parseBackendError = (value: unknown): BackendErrorPayload | null => {
  if (!isRecord(value)) {
    return null;
  }

  const { error, message } = value;
  if (typeof error !== "string" || typeof message !== "string") {
    return null;
  }

  return { error, message };
};

const buildUrl = (path: string): string => {
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  if (path.startsWith("/")) {
    return `${API_BASE}${path}`;
  }

  return `${API_BASE}/${path}`;
};

const shouldSetJsonContentType = (body: BodyInit | null | undefined): boolean => {
  if (!body) {
    return false;
  }

  if (typeof FormData !== "undefined" && body instanceof FormData) {
    return false;
  }

  return true;
};

export async function request<T>(path: string, init?: RequestOptions): Promise<T> {
  const { token, ...requestInit } = init ?? {};
  const headers = new Headers(requestInit.headers);

  if (token) {
    headers.set("Authorization", token);
  }

  if (shouldSetJsonContentType(requestInit.body) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  let response: Response;

  try {
    response = await fetch(buildUrl(path), {
      ...requestInit,
      headers,
    });
  } catch (error) {
    throw new ApiError(0, error instanceof Error ? error.message : "Ошибка сети");
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const contentType = response.headers.get("content-type") ?? "";
  const payload = contentType.includes("application/json")
    ? await response.json()
    : await response.text();

  if (!response.ok) {
    const backendError = parseBackendError(payload);
    if (backendError) {
      throw new ApiError(response.status, backendError.message, backendError.error);
    }

    throw new ApiError(
      response.status,
      typeof payload === "string" ? payload : `Запрос завершился со статусом ${response.status}`,
    );
  }

  return payload as T;
}
