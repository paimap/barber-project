'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface ProductFormProps {
  onSubmitSuccess: () => void;
}

export default function ProductForm({ onSubmitSuccess }: ProductFormProps) {
  const [loading, setLoading] = useState(false);
  const [displayPrice, setDisplayPrice] = useState(''); 
  const [rawPrice, setRawPrice] = useState(0);         

  // Fungsi format input ke Rupiah saat mengetik
  const handlePriceChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, ''); 
    const numericValue = parseInt(value) || 0;
    
    setRawPrice(numericValue);
    setDisplayPrice(new Intl.NumberFormat('id-ID').format(numericValue));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    
    const data = {
      name: formData.get('name'),
      price: rawPrice, 
    };

    try {
      const res = await fetch('/api/product', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert("Produk berhasil ditambahkan!");
        onSubmitSuccess();
      } else {
        const errData = await res.json();
        alert(`Error: ${errData.error || 'Gagal menyimpan'}`);
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan koneksi ke server.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.inputGroup}>
        <label>Nama Produk</label>
        <input 
          name="name" 
          type="text" 
          placeholder="e.g. Pomade" 
          required 
          className={styles.input}
        />
      </div>

      <div className={styles.inputGroup}>
        <label>Harga</label>
        <div className={styles.priceInputWrapper}>
          <span className={styles.currencyPrefix}>Rp</span>
          <input
            type="text"
            value={displayPrice}
            onChange={handlePriceChange}
            placeholder="10.000"
            required
            className={styles.inputWithPrefix}
          />
        </div>
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Menyimpan...' : 'Buat Produk'}
        </button>
      </div>
    </form>
  );
}