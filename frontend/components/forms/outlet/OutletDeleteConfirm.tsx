'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface OutletDeleteConfirmProps {
  id: number;
  address: string;
  onCancel: () => void;
  onSuccess: () => void;
}

export default function OutletDeleteConfirm({ id, address, onCancel, onSuccess }: OutletDeleteConfirmProps) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    setLoading(true);
    try {
      const res = await fetch(`/api/outlets/${id}`, {
        method: 'DELETE',
      });

      if (res.ok) {
        onSuccess();
      } else {
        const errData = await res.json();
        alert(errData.error || 'Gagal menghapus outlet');
      }
    } catch (err) {
      console.error('Delete error:', err);
      alert('Terjadi kesalahan koneksi');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.form}>
      <p>Hapus outlet <strong>{address}</strong>?</p>
      <p style={{ color: '#ef4444', fontSize: '13px', marginTop: '4px' }}>
        Tindakan ini akan menghapus outlet dan semua data terkait.
      </p>

      <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end', marginTop: '24px' }}>
        <button
          type="button"
          onClick={onCancel}
          className={styles.btnSecondary}
          style={{ padding: '8px 16px', borderRadius: '6px', border: '1px solid #e2e8f0', cursor: 'pointer', background: 'white' }}
          disabled={loading}
        >
          Batal
        </button>
        <button
          type="button"
          onClick={handleDelete}
          style={{ padding: '8px 16px', borderRadius: '6px', backgroundColor: '#ef4444', color: 'white', border: 'none', cursor: 'pointer' }}
          disabled={loading}
        >
          {loading ? 'Menghapus...' : 'Ya, Hapus'}
        </button>
      </div>
    </div>
  );
}
