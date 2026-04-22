export type UserRole = "buyer" | "seller" | "active" | "blocked" | string;

export interface User {
  id: string;
  fio: string;
  telephone_number: string;
  city: string;
  user_login: string;
  status: UserRole;
  user_description: string;
}

export interface UserResponseDto {
  user_id: string;
  fio: string;
  telephone_number: string;
  city: string;
  user_login: string;
  status: UserRole;
  user_description: string;
}

export interface UserUpsertPayload {
  fio: string;
  telephone_number: string;
  city: string;
  user_login: string;
  user_password: string;
  status: UserRole;
  user_description: string;
}

export interface LoginPayload {
  user_login: string;
  user_password: string;
}

export interface LoginResponse {
  token: string;
  user_id: string;
}

export interface RegisterResponse {
  id: string;
}
