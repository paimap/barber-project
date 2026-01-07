'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';
import { ProductOption } from '@/app/sales-barber/types';

interface ProductSalesFormProps {
  onSubmitSuccess: () => void;
  products: ProductOption[];
}

export default function ProductSalesForm({ onSubmitSuccess, products }: ProductSalesFormProps) {
  const [loading, setLoading] = useState(false);
  const [selectedItems, setSelectedItems] = useState<{ID: number, Quantity: number}[]>([]);

  const toggleProduct = (product: ProductOption) => {
    setSelectedItems(prev => {
      const exists = prev.find(item => item.ID === product.ID);
      if (exists) {
        return prev.filter(item => item.ID !== product.ID);
      }
      return [...prev, { ID: product.ID, Quantity: 1 }];
    });
  };

  const updateQuantity = (id: number, val: string) => {
    const qty = val === '' ? 0 : parseInt(val);
    
    if (isNaN(qty)) return;

    setSelectedItems(prev => prev.map(item => 
      item.ID === id ? { ...item, Quantity: qty } : item
    ));
  };

  const totalPrice = selectedItems.reduce((acc, item) => {
    const p = products.find(prod => prod.ID === item.ID);
    return acc + (p ? p.Price * item.Quantity : 0);
  }, 0);

  const isInvalid = selectedItems.length === 0 || selectedItems.some(item => item.Quantity <= 0);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isInvalid) return;

    setLoading(true);
    const formData = new FormData(e.currentTarget);
    
    const payload = {
      Products: selectedItems,
      PaymentType: formData.get('PaymentType'),
    };

    try {
      // Menggunakan Route Handler Next.js (app/api) untuk menghindari CORS
      const res = await fetch('/api/product-sales', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        alert('Penjualan Berhasil!');
        onSubmitSuccess();
      } else {
        const err = await res.json();
        alert(`Error: ${err.error || 'Gagal memproses penjualan'}`);
      }
    } catch (err) {
      alert('Koneksi ke server gagal.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.inputGroup}>
        <label className={styles.labelTitle}>Pilih Produk</label>
        <div className={styles.checkboxContainer}>
          {products?.map((prod) => {
            const itemInCart = selectedItems.find(item => item.ID === prod.ID);
            return (
              <div 
                key={prod.ID} 
                className={`${styles.serviceCard} ${itemInCart ? styles.selectedCard : ''}`}
                onClick={() => toggleProduct(prod)}
                style={{ justifyContent: 'space-between' }}
              >
                <div className={styles.serviceInfo}>
                  <span className={styles.serviceName}>
                    {itemInCart ? '✅ ' : ''}{prod.Name}
                  </span>
                  <span className={styles.servicePrice}>
                    Rp {prod.Price.toLocaleString('id-ID')}
                  </span>
                </div>

                {itemInCart && (
                  <div className={styles.qtyControl} onClick={(e) => e.stopPropagation()}>
                    <label style={{ fontSize: '10px', color: '#64748b' }}>Qty:</label>
                    <input 
                      type="number" 
                      min="1" 
                      // Mengubah 0 menjadi string kosong agar input bersih saat dihapus
                      value={itemInCart.Quantity === 0 ? '' : itemInCart.Quantity} 
                      onChange={(e) => updateQuantity(prod.ID, e.target.value)}
                      className={styles.inputQty}
                      placeholder="0"
                    />
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>

      <div className={styles.inputGroup}>
        <label className={styles.labelTitle}>Metode Pembayaran</label>
        <div className={styles.selectWrapper}>
          <select name="PaymentType" className={styles.selectInput} required defaultValue="CASH">
            <option value="CASH">💵 CASH</option>
            <option value="QRIS">📱 QRIS</option>
          </select>
        </div>
      </div>

      <div className={styles.totalSection}>
        <span className={styles.totalLabel}>Total Bayar</span>
        <span className={styles.totalAmount}>Rp {totalPrice.toLocaleString('id-ID')}</span>
      </div>

      <div className={styles.footer}>
        <button 
          type="submit" 
          disabled={loading || isInvalid} 
          className={styles.btnPrimary}
        >
          {loading ? 'Memproses...' : 'Konfirmasi Penjualan'}
        </button>
      </div>
    </form>
  );
}