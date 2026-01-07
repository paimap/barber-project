'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';
import { ServiceType } from '@/app/services-barber/types';

interface ServiceUpdateFormProps {
  initialData: ServiceType;
  onSubmitSuccess: () => void;
}

export default function ServiceUpdateForm({ initialData, onSubmitSuccess }: ServiceUpdateFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const payload = {
      PaymentType: formData.get('PaymentType'),
    };

    try {
      // Sesuaikan endpoint dengan API update Anda, biasanya menggunakan PUT atau PATCH
      const res = await fetch(`/api/services-barber/${initialData.ServiceId}`, {
        method: 'PUT', 
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        alert('Metode pembayaran berhasil diperbarui!');
        onSubmitSuccess();
      } else {
        const err = await res.json();
        alert(`Error: ${err.error || 'Gagal memperbarui data'}`);
      }
    } catch (err) {
      alert('Terjadi kesalahan koneksi.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.infoBox}>
        <p><strong>Service ID:</strong> #{initialData.ServiceId}</p>
        <p><strong>Layanan:</strong> {initialData.ServiceTypes.join(", ")}</p>
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
            <option value="CASH">💵 Tunai (CASH)</option>
            <option value="QRIS">📱 Non-Tunai (QRIS)</option>
          </select>
        </div>
      </div>

      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Menyimpan...' : 'Simpan Perubahan'}
        </button>
      </div>
    </form>
  );
}