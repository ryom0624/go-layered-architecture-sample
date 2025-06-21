import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Article, Category, Tag } from '../types';
import apiService from '../services/api';

const ArticleList: React.FC = () => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState('desc');

  useEffect(() => {
    loadData();
  }, []);

  useEffect(() => {
    searchArticles();
  }, [searchQuery, selectedCategory, selectedTags, sortBy, sortOrder]);

  const loadData = async () => {
    try {
      const [articlesRes, categoriesRes, tagsRes] = await Promise.all([
        apiService.getPublishedArticles(),
        apiService.getCategories(),
        apiService.getTags(),
      ]);
      
      setArticles(articlesRes.data);
      setCategories(categoriesRes.data);
      setTags(tagsRes.data);
    } catch (err: any) {
      setError('データの読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  const searchArticles = async () => {
    try {
      const params = {
        query: searchQuery || undefined,
        category_id: selectedCategory ? parseInt(selectedCategory) : undefined,
        tags: selectedTags.length > 0 ? selectedTags.join(',') : undefined,
        sort_by: sortBy,
        sort_order: sortOrder,
        status: 'published',
      };

      const response = await apiService.searchArticles(params);
      setArticles(response.data.data);
    } catch (err: any) {
      setError('検索に失敗しました');
    }
  };

  const toggleTag = (tagName: string) => {
    setSelectedTags(prev => 
      prev.includes(tagName) 
        ? prev.filter(t => t !== tagName)
        : [...prev, tagName]
    );
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP');
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">記事一覧</h1>
        
        {/* 検索フォーム */}
        <div className="bg-white p-6 rounded-lg shadow-md mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">検索</label>
              <input
                type="text"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="タイトルや内容で検索..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">カテゴリ</label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                value={selectedCategory}
                onChange={(e) => setSelectedCategory(e.target.value)}
              >
                <option value="">すべてのカテゴリ</option>
                {categories.map(category => (
                  <option key={category.id} value={category.id}>
                    {category.name}
                  </option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">並び順</label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                value={`${sortBy}-${sortOrder}`}
                onChange={(e) => {
                  const [field, order] = e.target.value.split('-');
                  setSortBy(field);
                  setSortOrder(order);
                }}
              >
                <option value="created_at-desc">新着順</option>
                <option value="created_at-asc">古い順</option>
                <option value="title-asc">タイトル順</option>
                <option value="author-asc">著者順</option>
              </select>
            </div>
          </div>
          
          {/* タグフィルタ */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">タグ</label>
            <div className="flex flex-wrap gap-2">
              {tags.map(tag => (
                <button
                  key={tag.id}
                  onClick={() => toggleTag(tag.name)}
                  className={`px-3 py-1 rounded-full text-sm ${
                    selectedTags.includes(tag.name)
                      ? 'bg-indigo-600 text-white'
                      : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                  }`}
                  style={selectedTags.includes(tag.name) ? {} : { backgroundColor: tag.color || '#e5e7eb' }}
                >
                  {tag.name}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* 記事リスト */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {articles.map(article => (
          <div key={article.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
            <div className="p-6">
              <div className="flex items-center justify-between mb-2">
                {article.category && (
                  <span className="inline-block bg-blue-100 text-blue-800 text-sm px-2 py-1 rounded">
                    {article.category.name}
                  </span>
                )}
                <span className="text-sm text-gray-500">
                  {formatDate(article.created_at)}
                </span>
              </div>
              
              <h2 className="text-xl font-semibold text-gray-900 mb-2">
                <Link 
                  to={`/articles/${article.id}`}
                  className="hover:text-indigo-600 transition-colors"
                >
                  {article.title}
                </Link>
              </h2>
              
              <p className="text-gray-600 mb-4 line-clamp-3">
                {article.content.substring(0, 150)}...
              </p>
              
              {article.author && (
                <div className="flex items-center justify-between">
                  <span className="text-sm text-gray-500">
                    著者: {article.author.name}
                  </span>
                  <div className="flex items-center space-x-2 text-sm text-gray-500">
                    {article.view_count && (
                      <span>👁 {article.view_count}</span>
                    )}
                    {article.favorite_count && (
                      <span>❤️ {article.favorite_count}</span>
                    )}
                  </div>
                </div>
              )}
              
              {article.tags && article.tags.length > 0 && (
                <div className="mt-3 flex flex-wrap gap-1">
                  {article.tags.map(tag => (
                    <span
                      key={tag.id}
                      className="inline-block text-xs px-2 py-1 rounded-full text-white"
                      style={{ backgroundColor: tag.color || '#6b7280' }}
                    >
                      {tag.name}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {articles.length === 0 && !isLoading && (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg">記事が見つかりませんでした</p>
        </div>
      )}
    </div>
  );
};

export default ArticleList;