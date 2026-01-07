'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';
import { ServiceTypeOption } from '@/app/services-barber/types';

interface ServiceFormProps {
  onSubmitSuccess: () => void;
  options: ServiceTypeOption[];
}

export default function ServiceForm({ onSubmitSuccess, options }: ServiceFormProps) {
  const [loading, setLoading] = useState(false);
  const [selectedTypes, setSelectedTypes] = useState<number[]>([]);

  // Fungsi toggle pilih/hapus layanan
  const handleTypeChange = (id: number) => {
    const targetId = Number(id);
    setSelectedTypes(prev => 
      prev.includes(targetId) 
        ? prev.filter(t => t !== targetId) 
        : [...prev, targetId]
    );
  };

  // Menghitung total harga otomatis berdasarkan state selectedTypes
  const totalPrice = selectedTypes.reduce((acc, currentId) => {
    const service = options.find(opt => opt.ID === currentId);
    return acc + (service ? Number(service.Price) : 0);
  }, 0);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (selectedTypes.length === 0) return alert("Pilih minimal satu layanan.");

    setLoading(true);
    const formData = new FormData(e.currentTarget);
    
    const payload = {
      ServiceTypeID: selectedTypes,
      PaymentType: formData.get('PaymentType'),
    };

    try {
      const res = await fetch('/api/services-barber', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        alert('Service Berhasil Dibuat!');
        onSubmitSuccess();
      } else {
        const err = await res.json();
        alert(`Error: ${err.error || 'Gagal membuat service'}`);
      }
    } catch (err) {
      alert('Koneksi ke server gagal.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      {/* SEKSI PILIH LAYANAN */}
      <div className={styles.inputGroup}>
        <label className={styles.labelTitle}>Pilih Layanan</label>
        <div className={styles.checkboxContainer}>
          {options && options.length > 0 ? (
            options.map((opt) => {
              const isChecked = selectedTypes.includes(opt.ID);
              return (
                <div 
                  key={opt.ID} 
                  className={`${styles.serviceCard} ${isChecked ? styles.selectedCard : ''}`}
                  onClick={() => handleTypeChange(opt.ID)}
                >
                  <div className={styles.serviceInfo}>
                    <span className={styles.serviceName}>
                        {opt.Name}
                    </span>
                    <span className={styles.servicePrice}>
                        Rp {Number(opt.Price).toLocaleString('id-ID')}
                    </span>
                  </div>
                </div>
              );
            })
          ) : (
            <p className={styles.emptyText}>Tidak ada tipe layanan tersedia.</p>
          )}
        </div>
      </div>

      {/* SEKSI METODE PEMBAYARAN */}
      <div className={styles.inputGroup}>
        <label className={styles.labelTitle}>Metode Pembayaran</label>
        <div className={styles.selectWrapper}>
          <select 
            name="PaymentType" 
            className={styles.selectInput} 
            required 
            defaultValue="CASH"
          >
            <option value="CASH">💵 Tunai (CASH)</option>
            <option value="QRIS">📱 Non-Tunai (QRIS)</option>
          </select>
        </div>
      </div>

      {/* SEKSI RINGKASAN HARGA */}
      <div className={styles.totalSection}>
        <span className={styles.totalLabel}>Total Transaksi</span>
        <span className={styles.totalAmount}>
          Rp {totalPrice.toLocaleString('id-ID')}
        </span>
      </div>

      <div className={styles.footer}>
        <button 
          type="submit" 
          disabled={loading || selectedTypes.length === 0} 
          className={styles.btnPrimary}
        >
          {loading ? 'Sedang Menyimpan...' : 'Konfirmasi & Buat Service'}
        </button>
      </div>
    </form>
  );
}