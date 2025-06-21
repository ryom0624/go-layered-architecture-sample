import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import Header from './components/Header';
import Login from './components/Login';
import Register from './components/Register';
import ArticleList from './components/ArticleList';
import ArticleDetail from './components/ArticleDetail';
import Dashboard from './components/Dashboard';
import ProtectedRoute from './components/ProtectedRoute';

const AppContent: React.FC = () => {
  const { isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<ArticleList />} />
          <Route path="/articles/:id" element={<ArticleDetail />} />
          <Route 
            path="/login" 
            element={
              isAuthenticated ? <Navigate to="/dashboard" replace /> : <Login />
            } 
          />
          <Route 
            path="/register" 
            element={
              isAuthenticated ? <Navigate to="/dashboard" replace /> : <Register />
            } 
          />
          <Route 
            path="/dashboard" 
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/favorites" 
            element={
              <ProtectedRoute>
                <div className="max-w-7xl mx-auto px-4 py-8">
                  <h1 className="text-3xl font-bold text-gray-900">お気に入り記事</h1>
                  <p className="mt-4 text-gray-600">この機能は実装中です。</p>
                </div>
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/reading-lists" 
            element={
              <ProtectedRoute>
                <div className="max-w-7xl mx-auto px-4 py-8">
                  <h1 className="text-3xl font-bold text-gray-900">読書リスト</h1>
                  <p className="mt-4 text-gray-600">この機能は実装中です。</p>
                </div>
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/categories" 
            element={
              <div className="max-w-7xl mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold text-gray-900">カテゴリ</h1>
                <p className="mt-4 text-gray-600">この機能は実装中です。</p>
              </div>
            } 
          />
          <Route 
            path="/tags" 
            element={
              <div className="max-w-7xl mx-auto px-4 py-8">
                <h1 className="text-3xl font-bold text-gray-900">タグ</h1>
                <p className="mt-4 text-gray-600">この機能は実装中です。</p>
              </div>
            } 
          />
          <Route 
            path="/profile" 
            element={
              <ProtectedRoute>
                <div className="max-w-7xl mx-auto px-4 py-8">
                  <h1 className="text-3xl font-bold text-gray-900">プロフィール</h1>
                  <p className="mt-4 text-gray-600">この機能は実装中です。</p>
                </div>
              </ProtectedRoute>
            } 
          />
          <Route 
            path="/settings" 
            element={
              <ProtectedRoute>
                <div className="max-w-7xl mx-auto px-4 py-8">
                  <h1 className="text-3xl font-bold text-gray-900">設定</h1>
                  <p className="mt-4 text-gray-600">この機能は実装中です。</p>
                </div>
              </ProtectedRoute>
            } 
          />
        </Routes>
      </main>
    </div>
  );
};

function App() {
  return (
    <AuthProvider>
      <Router>
        <AppContent />
      </Router>
    </AuthProvider>
  );
}

export default App;