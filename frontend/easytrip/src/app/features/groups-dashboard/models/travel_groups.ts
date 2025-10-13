
/**
 * Interface para os dados de criação de um novo grupo.
 * Corresponde ao TravelGroupCreateRequest no YAML.
 */
export interface TravelGroupCreateRequest {
  name: string;
  description?: string;
  start_date: string; // 'YYYY-MM-DD'
  end_date: string;   // 'YYYY-MM-DD'
}

/**
 * Interface para a resposta de listagem de grupos.
 * Corresponde ao TravelGroupListItem no YAML.
 */
export interface TravelGroupListItem {
  id: number;
  name: string;
  description: string;
  start_date: string; 
  end_date: string;   
  memberCount: number;
  creatorId: number;
  creatorName: string;
}

/**
 * Interface para a resposta de sucesso na criação de um grupo.
 * Corresponde ao TravelGroupDetails no YAML.
 */
export interface TravelGroupDetails {
  id: number;
  name: string;
  description: string;
  creatorId: number;
  start_date: string;
  end_date: string;
  createdAt: string;
}

export interface TravelGroupListItem {
  id: number;
  name: string;
  description: string;
  start_date: string; 
  end_date: string;   
  memberCount: number;
  creatorId: number;
  creatorName: string;
}
