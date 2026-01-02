'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface OutletFormProps {
  onSubmitSuccess: () => void;
}

export default function OutletForm({ onSubmitSuccess }: OutletFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const data = {
      address: formData.get('address'),
      phone: formData.get('phone'),
    };

    try {
      const res = await fetch('/api/outlets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert('Outlet berhasil dibuat!');
        onSubmitSuccess();
      } else {
        const errData = await res.json();
        alert(`Error: ${errData.error || 'Gagal membuat outlet'}`);
      }
    } catch (err) {
      console.error(err);
      alert('Terjadi kesalahan koneksi.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.inputGroup}>
        <label>Alamat</label>
        <input
          name="address"
          type="text"
          placeholder="Masukkan alamat outlet"
          required
          className={styles.input}
        />
      </div>

      <div className={styles.inputGroup}>
        <label>Nomor Telepon</label>
        <input
          name="phone"
          type="tel"
          placeholder="Masukkan nomor telepon"
          required
          className={styles.input}
        />
      </div>

      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Menyimpan...' : 'Buat Outlet'}
        </button>
      </div>
    </form>
  );
}
