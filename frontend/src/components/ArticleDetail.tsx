import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Article, Comment } from '../types';
import { useAuth } from '../contexts/AuthContext';
import apiService from '../services/api';

const ArticleDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { user, isAuthenticated } = useAuth();
  const [article, setArticle] = useState<Article | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [isFavorited, setIsFavorited] = useState(false);
  const [newComment, setNewComment] = useState('');
  const [replyTo, setReplyTo] = useState<number | null>(null);
  const [replyContent, setReplyContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (id) {
      loadArticle(parseInt(id));
      loadComments(parseInt(id));
      if (isAuthenticated) {
        checkFavoriteStatus(parseInt(id));
      }
    }
  }, [id, isAuthenticated]);

  const loadArticle = async (articleId: number) => {
    try {
      const response = await apiService.getArticle(articleId);
      setArticle(response.data);
    } catch (err: any) {
      setError('記事の読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  const loadComments = async (articleId: number) => {
    try {
      const response = await apiService.getArticleComments(articleId);
      setComments(response.data);
    } catch (err: any) {
      console.error('コメントの読み込みに失敗しました:', err);
    }
  };

  const checkFavoriteStatus = async (articleId: number) => {
    try {
      const response = await apiService.getFavoriteStatus(articleId);
      setIsFavorited(response.data.is_favorited);
    } catch (err: any) {
      console.error('お気に入り状態の確認に失敗しました:', err);
    }
  };

  const toggleFavorite = async () => {
    if (!article || !isAuthenticated) return;

    try {
      if (isFavorited) {
        await apiService.removeFromFavorites(article.id);
        setIsFavorited(false);
        if (article.favorite_count) {
          setArticle({ ...article, favorite_count: article.favorite_count - 1 });
        }
      } else {
        await apiService.addToFavorites(article.id);
        setIsFavorited(true);
        setArticle({ ...article, favorite_count: (article.favorite_count || 0) + 1 });
      }
    } catch (err: any) {
      console.error('お気に入り操作に失敗しました:', err);
    }
  };

  const submitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!article || !isAuthenticated || !newComment.trim()) return;

    setIsSubmitting(true);
    try {
      const response = await apiService.createComment(article.id, {
        content: newComment,
        user_id: user!.id,
      });
      setComments([...comments, response.data]);
      setNewComment('');
    } catch (err: any) {
      console.error('コメントの投稿に失敗しました:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const submitReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!replyTo || !replyContent.trim() || !isAuthenticated) return;

    setIsSubmitting(true);
    try {
      await apiService.createReply(replyTo, {
        content: replyContent,
        user_id: user!.id,
      });
      
      // コメントリストを更新
      await loadComments(article!.id);
      setReplyTo(null);
      setReplyContent('');
    } catch (err: any) {
      console.error('返信の投稿に失敗しました:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const renderComments = (comments: Comment[], depth = 0) => {
    return comments.map(comment => (
      <div key={comment.id} className={`border-l-2 border-gray-200 ${depth > 0 ? 'ml-8' : ''} pl-4 py-4`}>
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center">
            <span className="font-medium text-gray-900">
              {comment.user?.name || 'Anonymous'}
            </span>
            <span className="mx-2 text-gray-400">•</span>
            <span className="text-sm text-gray-500">
              {formatDate(comment.created_at)}
            </span>
          </div>
          {comment.status === 'pending' && (
            <span className="text-xs bg-yellow-100 text-yellow-800 px-2 py-1 rounded">
              承認待ち
            </span>
          )}
        </div>
        
        <p className="text-gray-800 mb-3">{comment.content}</p>
        
        {isAuthenticated && depth < 2 && (
          <button
            onClick={() => setReplyTo(comment.id)}
            className="text-sm text-indigo-600 hover:text-indigo-800"
          >
            返信
          </button>
        )}
        
        {replyTo === comment.id && (
          <form onSubmit={submitReply} className="mt-3">
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              rows={3}
              placeholder="返信を入力..."
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              required
            />
            <div className="mt-2 flex gap-2">
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-50"
              >
                {isSubmitting ? '投稿中...' : '返信'}
              </button>
              <button
                type="button"
                onClick={() => {
                  setReplyTo(null);
                  setReplyContent('');
                }}
                className="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
              >
                キャンセル
              </button>
            </div>
          </form>
        )}
        
        {comment.replies && comment.replies.length > 0 && (
          <div className="mt-4">
            {renderComments(comment.replies, depth + 1)}
          </div>
        )}
      </div>
    ));
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  if (error || !article) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error || '記事が見つかりません'}
        </div>
        <Link to="/" className="inline-block mt-4 text-indigo-600 hover:text-indigo-800">
          ← 記事一覧に戻る
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <Link to="/" className="inline-block mb-6 text-indigo-600 hover:text-indigo-800">
        ← 記事一覧に戻る
      </Link>
      
      <article className="bg-white rounded-lg shadow-md overflow-hidden">
        <div className="p-8">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-4">
              {article.category && (
                <span className="inline-block bg-blue-100 text-blue-800 px-3 py-1 rounded">
                  {article.category.name}
                </span>
              )}
              <span className="text-gray-500">
                {formatDate(article.created_at)}
              </span>
            </div>
            
            {isAuthenticated && (
              <button
                onClick={toggleFavorite}
                className={`flex items-center space-x-1 px-3 py-1 rounded ${
                  isFavorited 
                    ? 'bg-red-100 text-red-600' 
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                <span>{isFavorited ? '❤️' : '🤍'}</span>
                <span>{article.favorite_count || 0}</span>
              </button>
            )}
          </div>
          
          <h1 className="text-3xl font-bold text-gray-900 mb-4">{article.title}</h1>
          
          {article.author && (
            <div className="flex items-center mb-6">
              <span className="text-gray-700">著者: {article.author.name}</span>
              {article.view_count && (
                <span className="ml-4 text-gray-500">👁 {article.view_count} 回閲覧</span>
              )}
            </div>
          )}
          
          <div className="prose max-w-none mb-6">
            {article.content.split('\n').map((paragraph, index) => (
              <p key={index} className="mb-4 text-gray-800 leading-relaxed">
                {paragraph}
              </p>
            ))}
          </div>
          
          {article.tags && article.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mb-6">
              {article.tags.map(tag => (
                <span
                  key={tag.id}
                  className="inline-block px-3 py-1 rounded-full text-sm text-white"
                  style={{ backgroundColor: tag.color || '#6b7280' }}
                >
                  {tag.name}
                </span>
              ))}
            </div>
          )}
        </div>
      </article>
      
      {/* コメントセクション */}
      <div className="mt-8 bg-white rounded-lg shadow-md p-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-6">
          コメント ({comments.length})
        </h2>
        
        {isAuthenticated ? (
          <form onSubmit={submitComment} className="mb-8">
            <textarea
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              rows={4}
              placeholder="コメントを入力..."
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              required
            />
            <button
              type="submit"
              disabled={isSubmitting}
              className="mt-3 px-6 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 disabled:opacity-50"
            >
              {isSubmitting ? '投稿中...' : 'コメント投稿'}
            </button>
          </form>
        ) : (
          <div className="mb-8 p-4 bg-gray-50 rounded-lg">
            <p className="text-gray-600">
              コメントを投稿するには
              <Link to="/login" className="text-indigo-600 hover:text-indigo-800 mx-1">
                ログイン
              </Link>
              が必要です。
            </p>
          </div>
        )}
        
        <div className="space-y-4">
          {comments.length > 0 ? (
            renderComments(comments)
          ) : (
            <p className="text-gray-500 text-center py-8">
              まだコメントがありません。最初のコメントを投稿しませんか？
            </p>
          )}
        </div>
      </div>
    </div>
  );
};

export default ArticleDetail;