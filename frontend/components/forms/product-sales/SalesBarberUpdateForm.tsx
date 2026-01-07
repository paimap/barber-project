'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';
import { SalesBarber } from '@/app/sales-barber/types';

interface SalesBarberUpdateFormProps {
  initialData: SalesBarber;
  onSubmitSuccess: () => void;
}

export default function SalesBarberUpdateForm({ initialData, onSubmitSuccess }: SalesBarberUpdateFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const payload = {
      PaymentType: formData.get('PaymentType'),
    };

    try {
      // Pastikan endpoint ini sesuai dengan backend Go Anda untuk update sales
      const res = await fetch(`/api/product-sales/${initialData.ServiceId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        alert('Metode pembayaran berhasil diperbarui!');
        onSubmitSuccess();
      } else {
        const err = await res.json();
        alert(`Error: ${err.error || 'Gagal update'}`);
      }
    } catch (err) {
      alert('Koneksi gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.infoBox}>
        <p><strong>Sale ID:</strong> #{initialData.ServiceId}</p>
        <p><strong>Products:</strong> {initialData.ProductName.join(", ")}</p>
        <p><strong>Total:</strong> Rp {initialData.PriceAtSale.toLocaleString('id-ID')}</p>
      </div>

      <div className={styles.inputGroup}>
        <label className={styles.labelTitle}>Ubah Metode Pembayaran</label>
        <div className={styles.selectWrapper}>
          <select 
            name="PaymentType" 
            className={styles.selectInput} 
            required 
            defaultValue={initialData.PaymentType}
          >
            <option value="CASH">💵 CASH</option>
            <option value="QRIS">📱 QRIS</option>
          </select>
        </div>
      </div>

      <div className={styles.footer}>
        <button type="submit" disabled={loading} className={styles.btnPrimary}>
          {loading ? 'Menyimpan...' : 'Update Pembayaran'}
        </button>
      </div>
    </form>
  );
}