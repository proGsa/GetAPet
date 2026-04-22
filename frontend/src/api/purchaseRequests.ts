import { request } from "./client";
import type {
  PurchaseRequest,
  PurchaseRequestCreatePayload,
  PurchaseRequestStatusPayload,
} from "../types/purchaseRequest";

export const purchaseRequestsApi = {
  list: (token: string) =>
    request<PurchaseRequest[]>("/purchase-requests", {
      token,
    }),

  listByBuyer: (buyerId: string, token: string) =>
    request<PurchaseRequest[]>(`/purchase-requests/buyer/${buyerId}`, {
      token,
    }),

  listBySeller: (sellerId: string, token: string) =>
    request<PurchaseRequest[]>(`/purchase-requests/seller/${sellerId}`, {
      token,
    }),

  listByPet: (petId: string, token: string) =>
    request<PurchaseRequest[]>(`/purchase-requests/pet/${petId}`, {
      token,
    }),

  getById: (id: string, token: string) =>
    request<PurchaseRequest>(`/purchase-requests/${id}`, {
      token,
    }),

  create: (payload: PurchaseRequestCreatePayload, token: string) =>
    request<PurchaseRequest>("/purchase-requests", {
      method: "POST",
      body: JSON.stringify(payload),
      token,
    }),

  updateStatus: (id: string, payload: PurchaseRequestStatusPayload, token: string) =>
    request<PurchaseRequest>(`/purchase-requests/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify(payload),
      token,
    }),

  remove: (id: string, token: string) =>
    request<void>(`/purchase-requests/${id}`, {
      method: "DELETE",
      token,
    }),
};
