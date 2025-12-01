'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { isAuthenticated } from '@/lib/auth';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    if (isAuthenticated()) {
      router.push('/todos');
    } else {
      router.push('/login');
    }
  }, [router]);

  return (
    <div className="loading">
      <p>読み込み中...</p>
    </div>
  );
}
