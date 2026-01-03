'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';
import { ProductType } from '@/app/stock-outlet/types';


interface StockFormProps {
  products: ProductType[];
  onSubmitSuccess: () => void;
}

export default function StockForm({ products, onSubmitSuccess }: StockFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const data = {
      product_id: parseInt(formData.get('product_id') as string),
      quantity: parseInt(formData.get('quantity') as string),
    };

    try {
      const res = await fetch('/api/stock-outlet', { 
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert("Stok berhasil ditambahkan!");
        onSubmitSuccess();
      } else {
        const errData = await res.json();
        alert(`Error: ${errData.error || 'Gagal menyimpan stok'}`);
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan koneksi.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.inputGroup}>
            <label>Pilih Produk</label>
        <select 
            name="product_id" 
            required 
            className={styles.input}
            defaultValue="" 
        >
            <option value="" disabled>-- Pilih Produk --</option>
            {products.map((product) => (
                <option key={product.ID} value={product.ID}>
                {product.Name} — {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(product.Price)}
                </option>
            ))}
        </select>
        </div>

      <div className={styles.inputGroup}>
        <label>Jumlah Stok Masuk</label>
        <input 
          name="quantity" 
          type="number" 
          min="1"
          placeholder="e.g. 10" 
          required 
          className={styles.input}
        />
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Menyimpan...' : 'Tambah ke Inventori'}
        </button>
      </div>
    </form>
  );
}