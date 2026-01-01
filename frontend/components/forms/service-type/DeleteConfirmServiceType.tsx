'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface DeleteConfirmProps {
  id: string | number;
  name: string;
  onSuccess: () => void;
  onCancel: () => void;
}

export default function DeleteConfirmServiceType({ id, name, onSuccess, onCancel }: DeleteConfirmProps) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    setLoading(true);
    try {
      // Endpoint diganti ke service-type
      const res = await fetch(`/api/service-type/${id}`, {
        method: 'DELETE',
      });

      if (res.ok) {
        onSuccess();
      } else {
        const errData = await res.json();
        alert(errData.error || "Gagal menghapus layanan");
      }
    } catch (err) {
      console.error("Delete error:", err);
      alert("Terjadi kesalahan koneksi");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.form}>
      <p>Apakah Anda yakin ingin menghapus layanan <strong>{name}</strong>?</p>
      <p style={{ color: '#ef4444', fontSize: '13px', marginTop: '4px' }}>
        Tindakan ini tidak dapat dibatalkan.
      </p>
      
      <div style={{ display: 'flex', gap: '10px', justifyContent: 'flex-end', marginTop: '24px' }}>
        <button 
          type="button" 
          onClick={onCancel} 
          className={styles.btnSecondary}
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