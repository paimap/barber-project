'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface StockUpdateFormProps {
  initialData: {
    ProductID: number;
    ProductName: string;
    Quantity: number;
  };
  onSubmitSuccess: () => void;
}

export default function StockUpdateForm({ initialData, onSubmitSuccess }: StockUpdateFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const data = {
      quantity: parseInt(formData.get('quantity') as string),
    };

    try {
      const res = await fetch(`/api/stock/${initialData.ProductID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert("Jumlah stok berhasil diperbarui!");
        onSubmitSuccess();
      } else {
        const errData = await res.json();
        alert(errData.error || "Gagal memperbarui stok");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan jaringan");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.inputGroup}>
            <label>Nama Produk</label>
            <input 
                type="text" 
                value={initialData.ProductName} 
                readOnly // Lebih baik gunakan readOnly daripada disabled untuk form submission jika perlu
                className={styles.input}
                style={{ 
                backgroundColor: '#f8fafc', 
                color: '#64748b', 
                cursor: 'not-allowed',
                borderColor: '#e2e8f0' 
                }}
            />
        </div>

      <div className={styles.inputGroup}>
        <label>Jumlah Stok Saat Ini</label>
        <input 
          name="quantity" 
          type="number" 
          defaultValue={initialData.Quantity} 
          required 
          className={styles.input}
        />
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Updating...' : 'Simpan Perubahan'}
        </button>
      </div>
    </form>
  );
}