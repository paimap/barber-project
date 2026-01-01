'use client'

import { useState, useEffect } from 'react';
import styles from '@/components/forms/Form.module.css';

interface ProductUpdateFormProps {
  initialData: {
    ID: any;
    Name: string;
    Price: number; // Gunakan number agar konsisten dengan DB
  };
  onSubmitSuccess: () => void;
}

export default function ProductUpdateForm({ initialData, onSubmitSuccess }: ProductUpdateFormProps) {
  const [loading, setLoading] = useState(false);
  
  // State untuk menangani format tampilan harga
  const [displayPrice, setDisplayPrice] = useState('');
  const [rawPrice, setRawPrice] = useState(0);

  // Inisialisasi data awal ke dalam state
  useEffect(() => {
    if (initialData) {
      setRawPrice(initialData.Price);
      setDisplayPrice(new Intl.NumberFormat('id-ID').format(initialData.Price));
    }
  }, [initialData]);

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
    
    // Pastikan data yang dikirim murni angka untuk price
    const data = {
      name: formData.get('name'),
      price: rawPrice,
    };

    try {
      const res = await fetch(`/api/product/${initialData.ID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert("Produk berhasil diperbarui!");
        onSubmitSuccess();
      } else {
        const errData = await res.json();
        alert(errData.error || "Gagal memperbarui produk");
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
        <label>Name</label>
        <input 
          name="name" 
          type="text" 
          defaultValue={initialData.Name} 
          required 
          className={styles.input}
        />
      </div>

      <div className={styles.inputGroup}>
        <label>Price</label>
        <div className={styles.priceInputWrapper}>
          <span className={styles.currencyPrefix}>Rp</span>
          <input 
            type="text" 
            value={displayPrice} 
            onChange={handlePriceChange}
            required 
            className={styles.inputWithPrefix}
          />
        </div>
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Updating...' : 'Save Changes'}
        </button>
      </div>
    </form>
  );
}