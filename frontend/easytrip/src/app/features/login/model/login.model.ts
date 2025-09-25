export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

/** Modelos de Resposta */

export interface LoginResponse {
  token: string;
  // Você pode adicionar mais dados aqui, como user_id, name, etc.
}