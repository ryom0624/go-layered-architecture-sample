import axios from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type { 
  User, 
  Article, 
  Category, 
  Tag, 
  Comment, 
  AuthResponse, 
  LoginRequest, 
  RegisterRequest,
  ReadingList,
  PaginatedResponse,
  SearchParams
} from '../types';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: 'http://localhost:8080/api/v1',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem('access_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.api.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401) {
          const refreshToken = localStorage.getItem('refresh_token');
          if (refreshToken) {
            try {
              const response = await this.refreshToken(refreshToken);
              localStorage.setItem('access_token', response.data.access_token);
              error.config.headers.Authorization = `Bearer ${response.data.access_token}`;
              return this.api.request(error.config);
            } catch (refreshError) {
              localStorage.removeItem('access_token');
              localStorage.removeItem('refresh_token');
              window.location.href = '/login';
            }
          }
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(data: LoginRequest): Promise<AxiosResponse<AuthResponse>> {
    return this.api.post('/auth/login', data);
  }

  async register(data: RegisterRequest): Promise<AxiosResponse<AuthResponse>> {
    return this.api.post('/auth/register', data);
  }

  async refreshToken(refreshToken: string): Promise<AxiosResponse<{ access_token: string }>> {
    return this.api.post('/auth/refresh', { refresh_token: refreshToken });
  }

  async logout(): Promise<AxiosResponse<void>> {
    const refreshToken = localStorage.getItem('refresh_token');
    return this.api.post('/auth/logout', { refresh_token: refreshToken });
  }

  async logoutAll(): Promise<AxiosResponse<void>> {
    return this.api.post('/auth/logout-all');
  }

  // User endpoints
  async getUsers(): Promise<AxiosResponse<User[]>> {
    return this.api.get('/users');
  }

  async getUser(id: number): Promise<AxiosResponse<User>> {
    return this.api.get(`/users/${id}`);
  }

  async createUser(data: Partial<User>): Promise<AxiosResponse<User>> {
    return this.api.post('/users', data);
  }

  async updateUser(id: number, data: Partial<User>): Promise<AxiosResponse<User>> {
    return this.api.put(`/users/${id}`, data);
  }

  async deleteUser(id: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/users/${id}`);
  }

  // Article endpoints
  async getArticles(): Promise<AxiosResponse<Article[]>> {
    return this.api.get('/articles');
  }

  async getPublishedArticles(): Promise<AxiosResponse<Article[]>> {
    return this.api.get('/articles/published');
  }

  async getPopularArticles(): Promise<AxiosResponse<Article[]>> {
    return this.api.get('/articles/popular');
  }

  async getRecentArticles(): Promise<AxiosResponse<Article[]>> {
    return this.api.get('/articles/recent');
  }

  async getArticle(id: number): Promise<AxiosResponse<Article>> {
    return this.api.get(`/articles/${id}`);
  }

  async createArticle(data: Partial<Article>): Promise<AxiosResponse<Article>> {
    return this.api.post('/articles', data);
  }

  async updateArticle(id: number, data: Partial<Article>): Promise<AxiosResponse<Article>> {
    return this.api.put(`/articles/${id}`, data);
  }

  async deleteArticle(id: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/articles/${id}`);
  }

  async publishArticle(id: number): Promise<AxiosResponse<Article>> {
    return this.api.put(`/articles/${id}/publish`);
  }

  async unpublishArticle(id: number): Promise<AxiosResponse<Article>> {
    return this.api.put(`/articles/${id}/unpublish`);
  }

  // Search endpoint
  async searchArticles(params: SearchParams): Promise<AxiosResponse<PaginatedResponse<Article>>> {
    return this.api.get('/search', { params });
  }

  // Comment endpoints
  async getArticleComments(articleId: number): Promise<AxiosResponse<Comment[]>> {
    return this.api.get(`/articles/${articleId}/comments`);
  }

  async createComment(articleId: number, data: Partial<Comment>): Promise<AxiosResponse<Comment>> {
    return this.api.post(`/articles/${articleId}/comments`, data);
  }

  async updateComment(id: number, data: Partial<Comment>): Promise<AxiosResponse<Comment>> {
    return this.api.put(`/comments/${id}`, data);
  }

  async deleteComment(id: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/comments/${id}`);
  }

  async createReply(commentId: number, data: Partial<Comment>): Promise<AxiosResponse<Comment>> {
    return this.api.post(`/comments/${commentId}/replies`, data);
  }

  // Category endpoints
  async getCategories(): Promise<AxiosResponse<Category[]>> {
    return this.api.get('/categories');
  }

  async getCategoriesWithCount(): Promise<AxiosResponse<(Category & { article_count: number })[]>> {
    return this.api.get('/categories/with-count');
  }

  async getCategory(slug: string): Promise<AxiosResponse<Category>> {
    return this.api.get(`/categories/${slug}`);
  }

  async createCategory(data: Partial<Category>): Promise<AxiosResponse<Category>> {
    return this.api.post('/categories', data);
  }

  async updateCategory(id: number, data: Partial<Category>): Promise<AxiosResponse<Category>> {
    return this.api.put(`/categories/${id}`, data);
  }

  async deleteCategory(id: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/categories/${id}`);
  }

  // Tag endpoints
  async getTags(): Promise<AxiosResponse<Tag[]>> {
    return this.api.get('/tags');
  }

  async getPopularTags(): Promise<AxiosResponse<(Tag & { usage_count: number })[]>> {
    return this.api.get('/tags/popular');
  }

  async getTag(slug: string): Promise<AxiosResponse<Tag>> {
    return this.api.get(`/tags/${slug}`);
  }

  async createTag(data: Partial<Tag>): Promise<AxiosResponse<Tag>> {
    return this.api.post('/tags', data);
  }

  async updateTag(id: number, data: Partial<Tag>): Promise<AxiosResponse<Tag>> {
    return this.api.put(`/tags/${id}`, data);
  }

  async deleteTag(id: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/tags/${id}`);
  }

  // Favorite endpoints
  async addToFavorites(articleId: number): Promise<AxiosResponse<void>> {
    return this.api.post(`/articles/${articleId}/favorite`);
  }

  async removeFromFavorites(articleId: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/articles/${articleId}/favorite`);
  }

  async getUserFavorites(userId: number): Promise<AxiosResponse<Article[]>> {
    return this.api.get(`/users/${userId}/favorites`);
  }

  async getFavoriteStatus(articleId: number): Promise<AxiosResponse<{ is_favorited: boolean }>> {
    return this.api.get(`/articles/${articleId}/favorite-status`);
  }

  // Reading List endpoints
  async getReadingLists(): Promise<AxiosResponse<ReadingList[]>> {
    return this.api.get('/reading-lists');
  }

  async getReadingList(id: number): Promise<AxiosResponse<ReadingList>> {
    return this.api.get(`/reading-lists/${id}`);
  }

  async createReadingList(data: Partial<ReadingList>): Promise<AxiosResponse<ReadingList>> {
    return this.api.post('/reading-lists', data);
  }

  async updateReadingList(id: number, data: Partial<ReadingList>): Promise<AxiosResponse<ReadingList>> {
    return this.api.put(`/reading-lists/${id}`, data);
  }

  async deleteReadingList(id: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/reading-lists/${id}`);
  }

  async addArticleToReadingList(listId: number, articleId: number): Promise<AxiosResponse<void>> {
    return this.api.post(`/reading-lists/${listId}/articles`, { article_id: articleId });
  }

  async removeArticleFromReadingList(listId: number, articleId: number): Promise<AxiosResponse<void>> {
    return this.api.delete(`/reading-lists/${listId}/articles/${articleId}`);
  }

  async getReadingListArticles(listId: number): Promise<AxiosResponse<Article[]>> {
    return this.api.get(`/reading-lists/${listId}/articles`);
  }
}

export default new ApiService();