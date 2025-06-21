export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface Article {
  id: number;
  title: string;
  content: string;
  author_id: number;
  author?: User;
  category_id?: number;
  category?: Category;
  tags?: Tag[];
  status: string;
  published_at?: string;
  created_at: string;
  updated_at: string;
  favorite_count?: number;
  view_count?: number;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
  color?: string;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  content: string;
  article_id: number;
  user_id: number;
  user?: User;
  parent_id?: number;
  replies?: Comment[];
  status: string;
  created_at: string;
  updated_at: string;
}

export interface AuthUser {
  id: number;
  name: string;
  email: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  user: AuthUser;
  access_token: string;
  refresh_token: string;
}

export interface ReadingList {
  id: number;
  name: string;
  description?: string;
  user_id: number;
  is_public: boolean;
  created_at: string;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface SearchParams {
  query?: string;
  author_id?: number;
  category_id?: number;
  tags?: string;
  status?: string;
  date_from?: string;
  date_to?: string;
  sort_by?: string;
  sort_order?: string;
  page?: number;
  limit?: number;
}