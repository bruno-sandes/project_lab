// Define a interface para os dados de login que serão enviados
export interface LoginRequest {
  email: string;
  password: string;
}

// Define a interface para a resposta de login da API
export interface LoginResponse {
  token: string;
}

// Define a interface para os dados de registro que serão enviados
export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}