'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { todoAPI } from '@/lib/api';
import { isAuthenticated, logout, getUser, isAdmin } from '@/lib/auth';
import { Todo } from '@/types';
import Link from 'next/link';

export default function Todos() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const router = useRouter();
  const user = getUser();

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }
    fetchTodos();
  }, [router]);

  const fetchTodos = async () => {
    try {
      const data = await todoAPI.getTodos();
      setTodos(data);
    } catch (err) {
      setError('TODOの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleAddTodo = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    try {
      await todoAPI.createTodo({
        title,
        description,
        completed: false,
      });
      setTitle('');
      setDescription('');
      fetchTodos();
    } catch (err) {
      setError('TODOの追加に失敗しました');
    }
  };

  const handleToggle = async (todo: Todo) => {
    try {
      await todoAPI.updateTodo(todo.id, {
        ...todo,
        completed: !todo.completed,
      });
      fetchTodos();
    } catch (err) {
      setError('TODOの更新に失敗しました');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('このTODOを削除してもよろしいですか？')) return;

    try {
      await todoAPI.deleteTodo(id);
      fetchTodos();
    } catch (err) {
      setError('TODOの削除に失敗しました');
    }
  };

  if (loading) {
    return <div className="loading">読み込み中...</div>;
  }

  return (
    <>
      <nav className="navbar">
        <div className="navbar-content">
          <h1>TODO管理</h1>
          <div className="navbar-menu">
            <span>{user?.email}</span>
            {isAdmin() && (
              <Link href="/admin">管理者ページ</Link>
            )}
            <button onClick={logout} className="btn btn-secondary">
              ログアウト
            </button>
          </div>
        </div>
      </nav>

      <div className="container">
        <div className="add-todo-form">
          <h2>新しいTODOを追加</h2>
          <form onSubmit={handleAddTodo}>
            <div className="form-group">
              <label htmlFor="title">タイトル</label>
              <input
                type="text"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="description">説明</label>
              <textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
            <button type="submit" className="btn btn-primary">
              追加
            </button>
          </form>
        </div>

        {error && <div className="error">{error}</div>}

        {todos.length === 0 ? (
          <div className="empty-state">
            <h3>TODOがありません</h3>
            <p>上のフォームから新しいTODOを追加してください</p>
          </div>
        ) : (
          <div className="todo-list">
            {todos.map((todo) => (
              <div key={todo.id} className={`todo-item ${todo.completed ? 'completed' : ''}`}>
                <input
                  type="checkbox"
                  className="todo-checkbox"
                  checked={todo.completed}
                  onChange={() => handleToggle(todo)}
                />
                <div className={`todo-content ${todo.completed ? 'completed' : ''}`}>
                  <h3>{todo.title}</h3>
                  {todo.description && <p>{todo.description}</p>}
                </div>
                <div className="todo-actions">
                  <button
                    onClick={() => handleDelete(todo.id)}
                    className="btn btn-danger"
                  >
                    削除
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </>
  );
}
