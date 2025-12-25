'use client'

import Table from '@/components/table/Table';
import styles from './Mitra.module.css'; // Kita buat CSS khusus layout page
import { Plus } from 'lucide-react';

export default function MitraClient() {
  const mitraHeaders = [
    { key: 'name', label: 'Name' },
    { key: 'phone', label: 'Phone Number' },
    { key: 'outlets', label: 'Outlet Count' },
    { key: 'joinDate', label: 'Join Date' },
  ];

  const mitraData = [
    { id: 1, name: "Budi Santoso", phone: "+62 812-3456-7890", outlets: "3 Outlets", joinDate: "12 Jan 2024" },
    { id: 2, name: "Andi Wijaya", phone: "+62 857-1122-3344", outlets: "1 Outlet", joinDate: "05 Feb 2024" },
  ];

  return (
    <div className={styles.pageContainer}>
      <div className={styles.headerSection}>
        <div className={styles.titleGroup}>
          <h1>Mitra List</h1>
          <p>Manage and monitor all your partner barbershops.</p>
        </div>
        <button className={styles.btnCreate} onClick={() => alert('Create Mitra')}>
          <Plus size={18} />
          <span>Create Mitra</span>
        </button>
      </div>

      <Table 
        headers={mitraHeaders} 
        data={mitraData}
        onDetail={(id) => console.log('View detail', id)}
        onUpdate={(id) => console.log('Update', id)}
        onDelete={(id) => console.log('Delete', id)}
      />
    </div>
  );
}