'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface ServiceTypeFormProps {
  onSubmitSuccess: () => void;
}

export default function ServiceTypeForm({ onSubmitSuccess }: ServiceTypeFormProps) {
  const [loading, setLoading] = useState(false);
  const [displayPrice, setDisplayPrice] = useState(''); 
  const [rawPrice, setRawPrice] = useState(0);         

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
      // Endpoint diganti ke service-type
      const res = await fetch('/api/service-type', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert("Layanan baru berhasil ditambahkan!");
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
        <label>Nama Layanan</label>
        <input 
          name="name" 
          type="text" 
          placeholder="e.g. Haircut Premium" 
          required 
          className={styles.input}
        />
      </div>

      <div className={styles.inputGroup}>
        <label>Harga Layanan</label>
        <div className={styles.priceInputWrapper}>
          <span className={styles.currencyPrefix}>Rp</span>
          <input
            type="text"
            value={displayPrice}
            onChange={handlePriceChange}
            placeholder="50.000"
            required
            className={styles.inputWithPrefix}
          />
        </div>
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Menyimpan...' : 'Tambah Layanan'}
        </button>
      </div>
    </form>
  );
}