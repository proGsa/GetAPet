export interface PurchaseRequest {
  id: string;
  pet_id: string;
  buyer_id: string;
  seller_id: string;
  status: string;
  request_date: string;
}

export interface PurchaseRequestCreatePayload {
  pet_id: string;
  buyer_id: string;
}

export interface PurchaseRequestStatusPayload {
  status: string;
}
