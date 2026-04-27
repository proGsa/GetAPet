import { request } from "./client";
import type {
  Pet,
  PetCreatePayload,
  PetCreateResponse,
  PetUpdatePayload,
  PetUpdateResponse,
} from "../types/pet";

export const petsApi = {
  list: () => request<Pet[]>("/pets"),

  getById: (id: string) => request<Pet>(`/pets/${id}`),

  create: (payload: PetCreatePayload, token: string) =>
    request<PetCreateResponse>("/pets", {
      method: "POST",
      body: JSON.stringify(payload),
      token,
    }),

  update: (id: string, payload: PetUpdatePayload, token: string) =>
    request<PetUpdateResponse>(`/pets/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
      token,
    }),

  remove: (id: string, token: string) =>
    request<void>(`/pets/${id}`, {
      method: "DELETE",
      token,
    }),
};
