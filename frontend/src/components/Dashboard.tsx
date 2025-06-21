import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Article, User } from '../types';
import { useAuth } from '../contexts/AuthContext';
import apiService from '../services/api';

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [recentArticles, setRecentArticles] = useState<Article[]>([]);
  const [popularArticles, setPopularArticles] = useState<Article[]>([]);
  const [favoriteArticles, setFavoriteArticles] = useState<Article[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [recentRes, popularRes, usersRes] = await Promise.all([
        apiService.getRecentArticles(),
        apiService.getPopularArticles(),
        apiService.getUsers(),
      ]);

      setRecentArticles(recentRes.data.slice(0, 5));
      setPopularArticles(popularRes.data.slice(0, 5));
      setUsers(usersRes.data);

      if (user) {
        try {
          const favoritesRes = await apiService.getUserFavorites(user.id);
          setFavoriteArticles(favoritesRes.data.slice(0, 5));
        } catch (err) {
          console.error('お気に入り記事の取得に失敗しました:', err);
        }
      }
    } catch (err: any) {
      setError('データの読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
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
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          ダッシュボード
        </h1>
        <p className="text-gray-600">
          ようこそ、{user?.name}さん！
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* 最新記事 */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900">最新記事</h2>
            <Link to="/" className="text-indigo-600 hover:text-indigo-800 text-sm">
              すべて見る
            </Link>
          </div>
          <div className="space-y-4">
            {recentArticles.map(article => (
              <div key={article.id} className="border-b border-gray-200 pb-4 last:border-b-0">
                <h3 className="font-medium text-gray-900 mb-1">
                  <Link 
                    to={`/articles/${article.id}`}
                    className="hover:text-indigo-600"
                  >
                    {article.title}
                  </Link>
                </h3>
                <div className="flex items-center text-sm text-gray-500">
                  <span>{article.author?.name}</span>
                  <span className="mx-2">•</span>
                  <span>{formatDate(article.created_at)}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 人気記事 */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900">人気記事</h2>
            <Link to="/articles/popular" className="text-indigo-600 hover:text-indigo-800 text-sm">
              すべて見る
            </Link>
          </div>
          <div className="space-y-4">
            {popularArticles.map(article => (
              <div key={article.id} className="border-b border-gray-200 pb-4 last:border-b-0">
                <h3 className="font-medium text-gray-900 mb-1">
                  <Link 
                    to={`/articles/${article.id}`}
                    className="hover:text-indigo-600"
                  >
                    {article.title}
                  </Link>
                </h3>
                <div className="flex items-center text-sm text-gray-500">
                  <span>{article.author?.name}</span>
                  <span className="mx-2">•</span>
                  {article.view_count && (
                    <>
                      <span>👁 {article.view_count}</span>
                      <span className="mx-2">•</span>
                    </>
                  )}
                  {article.favorite_count && (
                    <span>❤️ {article.favorite_count}</span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* お気に入り記事 */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900">お気に入り記事</h2>
            <Link to="/favorites" className="text-indigo-600 hover:text-indigo-800 text-sm">
              すべて見る
            </Link>
          </div>
          <div className="space-y-4">
            {favoriteArticles.length > 0 ? (
              favoriteArticles.map(article => (
                <div key={article.id} className="border-b border-gray-200 pb-4 last:border-b-0">
                  <h3 className="font-medium text-gray-900 mb-1">
                    <Link 
                      to={`/articles/${article.id}`}
                      className="hover:text-indigo-600"
                    >
                      {article.title}
                    </Link>
                  </h3>
                  <div className="flex items-center text-sm text-gray-500">
                    <span>{article.author?.name}</span>
                    <span className="mx-2">•</span>
                    <span>{formatDate(article.created_at)}</span>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-gray-500 text-sm">
                まだお気に入り記事がありません。
                <Link to="/" className="text-indigo-600 hover:text-indigo-800 ml-1">
                  記事を探す
                </Link>
              </p>
            )}
          </div>
        </div>

        {/* ユーザー統計 */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">プラットフォーム統計</h2>
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center">
              <div className="text-3xl font-bold text-indigo-600">
                {recentArticles.length + popularArticles.length}
              </div>
              <div className="text-sm text-gray-500">記事数</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-green-600">
                {users.length}
              </div>
              <div className="text-sm text-gray-500">ユーザー数</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-orange-600">
                {favoriteArticles.length}
              </div>
              <div className="text-sm text-gray-500">お気に入り</div>
            </div>
            <div className="text-center">
              <div className="text-3xl font-bold text-purple-600">
                {popularArticles.reduce((sum, article) => sum + (article.view_count || 0), 0)}
              </div>
              <div className="text-sm text-gray-500">総閲覧数</div>
            </div>
          </div>
        </div>
      </div>

      {/* クイックアクション */}
      <div className="mt-8 bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">クイックアクション</h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <Link 
            to="/reading-lists"
            className="block p-4 border border-gray-200 rounded-lg hover:border-indigo-300 hover:shadow-md transition-all"
          >
            <div className="text-center">
              <div className="text-2xl mb-2">📚</div>
              <div className="font-medium text-gray-900">読書リスト</div>
              <div className="text-sm text-gray-500">記事をリストに整理</div>
            </div>
          </Link>
          
          <Link 
            to="/search"
            className="block p-4 border border-gray-200 rounded-lg hover:border-indigo-300 hover:shadow-md transition-all"
          >
            <div className="text-center">
              <div className="text-2xl mb-2">🔍</div>
              <div className="font-medium text-gray-900">検索</div>
              <div className="text-sm text-gray-500">記事を探す</div>
            </div>
          </Link>
          
          <Link 
            to="/categories"
            className="block p-4 border border-gray-200 rounded-lg hover:border-indigo-300 hover:shadow-md transition-all"
          >
            <div className="text-center">
              <div className="text-2xl mb-2">📂</div>
              <div className="font-medium text-gray-900">カテゴリ</div>
              <div className="text-sm text-gray-500">分野別に閲覧</div>
            </div>
          </Link>
          
          <Link 
            to="/tags"
            className="block p-4 border border-gray-200 rounded-lg hover:border-indigo-300 hover:shadow-md transition-all"
          >
            <div className="text-center">
              <div className="text-2xl mb-2">🏷️</div>
              <div className="font-medium text-gray-900">タグ</div>
              <div className="text-sm text-gray-500">トピック別に閲覧</div>
            </div>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;