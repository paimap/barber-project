'use client'

import { useState, useEffect } from 'react';
import styles from '@/components/forms/Form.module.css';

interface ServiceTypeUpdateFormProps {
  initialData: {
    id: any;      // Menggunakan id (kecil) karena sudah dimapping di ServicesClient
    name: string;
    price: number;
  };
  onSubmitSuccess: () => void;
}

export default function ServiceTypeUpdateForm({ initialData, onSubmitSuccess }: ServiceTypeUpdateFormProps) {
  const [loading, setLoading] = useState(false);
  const [displayPrice, setDisplayPrice] = useState('');
  const [rawPrice, setRawPrice] = useState(0);

  useEffect(() => {
    if (initialData) {
      setRawPrice(initialData.price);
      setDisplayPrice(new Intl.NumberFormat('id-ID').format(initialData.price));
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
    const data = {
      name: formData.get('name'),
      price: rawPrice,
    };

    try {
      // Endpoint diganti ke service-type
      const res = await fetch(`/api/service-type/${initialData.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        alert("Informasi layanan berhasil diperbarui!");
        onSubmitSuccess();
      } else {
        const errData = await res.json();
        alert(errData.error || "Gagal memperbarui layanan");
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
        <label>Nama Layanan</label>
        <input 
          name="name" 
          type="text" 
          defaultValue={initialData.name} 
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
            required 
            className={styles.inputWithPrefix}
          />
        </div>
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Updating...' : 'Simpan Perubahan'}
        </button>
      </div>
    </form>
  );
}