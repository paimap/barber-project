'use client'

import Table from '@/components/table/Table';
import styles from './Services.module.css';
import { Plus } from 'lucide-react';

export default function ServicesClient() {
  const commonHeaders = [
    { key: 'name', label: 'Item Name' },
    { key: 'price', label: 'Price' },
  ];

  const services = [
    { id: 1, name: "Haircut Premium", price: "Rp 50.000" },
    { id: 2, name: "Shaving & Massage", price: "Rp 30.000" },
  ];

  const products = [
    { id: 101, name: "Pomade Waterbased", price: "Rp 85.000" },
    { id: 102, name: "Beard Oil", price: "Rp 120.000" },
  ];

  return (
    <div className={styles.pageContainer}>
      <header className={styles.headerSection}>
        <div className={styles.titleGroup}>
          <h1>Services & Products</h1>
          <p>Configure your barber offerings and shop inventory.</p>
        </div>
      </header>

      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Barber Services</h2>
          <button className={styles.btnAddSmall} onClick={() => alert('Add Service')}>
            <Plus size={16} /> Add Service
          </button>
        </div>
        <Table 
          headers={commonHeaders} 
          data={services}
          onUpdate={(id) => console.log('Update', id)}
          onDelete={(id) => console.log('Delete', id)}
        />
      </section>

      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Shop Products</h2>
          <button className={styles.btnAddSmall} onClick={() => alert('Add Product')}>
            <Plus size={16} /> Add Product
          </button>
        </div>
        <Table 
          headers={commonHeaders} 
          data={products}
          onUpdate={(id) => console.log('Update', id)}
          onDelete={(id) => console.log('Delete', id)}
        />
      </section>
    </div>
  );
}