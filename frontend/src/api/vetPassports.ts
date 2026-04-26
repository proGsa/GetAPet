import { request } from "./client";
import type {
  VetPassport,
  VetPassportCreateResponse,
  VetPassportUpsertPayload,
} from "../types/vetPassport";

export const vetPassportsApi = {
  list: () => request<VetPassport[]>("/vet-passports"),

  getById: (id: string) => request<VetPassport>(`/vet-passports/${id}`),

  create: (payload: VetPassportUpsertPayload) =>
    request<VetPassportCreateResponse>("/vet-passports", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  update: (id: string, payload: VetPassportUpsertPayload) =>
    request<VetPassport>(`/vet-passports/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    }),

  remove: (id: string) =>
    request<void>(`/vet-passports/${id}`, {
      method: "DELETE",
    }),
};
