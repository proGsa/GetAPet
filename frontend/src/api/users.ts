import { request } from "./client";
import type {
  LoginPayload,
  LoginResponse,
  RegisterResponse,
  User,
  UserResponseDto,
  UserUpsertPayload,
} from "../types/user";

const mapUser = (dto: UserResponseDto): User => ({
  id: dto.user_id,
  fio: dto.fio,
  telephone_number: dto.telephone_number,
  city: dto.city,
  user_login: dto.user_login,
  status: dto.status,
  user_description: dto.user_description,
});

export const usersApi = {
  register: (payload: UserUpsertPayload) =>
    request<RegisterResponse>("/users", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  login: (payload: LoginPayload) =>
    request<LoginResponse>("/users/login", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  list: async (token: string) => {
    const response = await request<UserResponseDto[]>("/users", { token });
    return response.map(mapUser);
  },

  getById: async (id: string, token: string) => {
    const response = await request<UserResponseDto>(`/users/${id}`, { token });
    return mapUser(response);
  },

  update: async (id: string, payload: UserUpsertPayload, token: string) => {
    const response = await request<UserResponseDto>(`/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
      token,
    });
    return mapUser(response);
  },

  remove: (id: string, token: string) =>
    request<void>(`/users/${id}`, {
      method: "DELETE",
      token,
    }),
};
