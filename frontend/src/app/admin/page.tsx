'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { adminAPI } from '@/lib/api';
import { isAuthenticated, isAdmin, logout, getUser } from '@/lib/auth';
import { User } from '@/types';
import Link from 'next/link';

export default function Admin() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const router = useRouter();
  const currentUser = getUser();

  useEffect(() => {
    if (!isAuthenticated()) {
      router.push('/login');
      return;
    }
    if (!isAdmin()) {
      router.push('/todos');
      return;
    }
    fetchUsers();
  }, [router]);

  const fetchUsers = async () => {
    try {
      const data = await adminAPI.getAllUsers();
      setUsers(data);
    } catch (err) {
      setError('ユーザーの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteUser = async (id: number) => {
    if (id === currentUser?.id) {
      alert('自分自身を削除することはできません');
      return;
    }

    if (!confirm('このユーザーを削除してもよろしいですか？')) return;

    try {
      await adminAPI.deleteUser(id);
      fetchUsers();
    } catch (err) {
      setError('ユーザーの削除に失敗しました');
    }
  };

  const handleToggleAdmin = async (user: User) => {
    if (user.id === currentUser?.id) {
      alert('自分自身の管理者権限は変更できません');
      return;
    }

    try {
      await adminAPI.updateUserRole(user.id, !user.is_admin);
      fetchUsers();
    } catch (err) {
      setError('権限の更新に失敗しました');
    }
  };

  if (loading) {
    return <div className="loading">読み込み中...</div>;
  }

  return (
    <>
      <nav className="navbar">
        <div className="navbar-content">
          <h1>管理者ページ</h1>
          <div className="navbar-menu">
            <Link href="/todos">TODOページ</Link>
            <span>{currentUser?.email}</span>
            <button onClick={logout} className="btn btn-secondary">
              ログアウト
            </button>
          </div>
        </div>
      </nav>

      <div className="container">
        <h2>ユーザー管理</h2>
        {error && <div className="error">{error}</div>}

        <div className="user-table">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>メールアドレス</th>
                <th>権限</th>
                <th>登録日</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>{user.email}</td>
                  <td>
                    <span className={`badge ${user.is_admin ? 'badge-admin' : 'badge-user'}`}>
                      {user.is_admin ? '管理者' : '一般ユーザー'}
                    </span>
                  </td>
                  <td>{new Date(user.created_at).toLocaleString('ja-JP')}</td>
                  <td>
                    <div className="action-buttons">
                      <button
                        onClick={() => handleToggleAdmin(user)}
                        className="btn btn-secondary"
                        disabled={user.id === currentUser?.id}
                      >
                        {user.is_admin ? '管理者解除' : '管理者に設定'}
                      </button>
                      <button
                        onClick={() => handleDeleteUser(user.id)}
                        className="btn btn-danger"
                        disabled={user.id === currentUser?.id}
                      >
                        削除
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </>
  );
}
