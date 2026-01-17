'use client'

import { useState, useEffect } from 'react';
import DeleteConfirm from '@/components/forms/mitra/DeleteConfirm';
import { useRouter } from 'next/navigation';
import Table from '@/components/table/Table';
import styles from './Mitra.module.css';
import { Plus } from 'lucide-react';
import { MitraClientProps } from './types';
import Modal from '@/components/modal/Modal';
import MitraForm from '@/components/forms/mitra/MitraForm';
import MitraUpdateForm from '@/components/forms/mitra/MitraUpdateForm';
import { StatCard } from '@/components/statcard/StatCard';
import { Wallet, Percent, ShoppingCart, TrendingUp } from 'lucide-react';

export default function MitraClient({ mitraData, summaryData }: MitraClientProps) {
  const [showModal, setShowModal] = useState(false);
  const router = useRouter();
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [selectedMitra, setSelectedMitra] = useState<{id: any, name: string} | null>(null);
  const [isUpdateOpen, setIsUpdateOpen] = useState(false);
  const [editData, setEditData] = useState<any>(null);

  const handleOpenUpdate = (id: any) => {
    const target = mitraData.find((m: any) => m.id === id);
    if (target) {
      setEditData(target);
      setIsUpdateOpen(true);
    }
  };

  const handleUpdateSuccess = () => {
    setIsUpdateOpen(false);
    setEditData(null);
    router.refresh();
  };

  const handleRefresh = () => {
    setShowModal(false);
    router.refresh(); 
  };

  const handleOpenDelete = (id: any) => {
    const mitra = mitraData.find((m) => m.id === id);
    if (mitra) {
      setSelectedMitra({ id: mitra.id, name: mitra.name });
      setIsDeleteOpen(true);
    }
  };

  const handleOpenView = (id: any) => {
    const mitra = mitraData.find((m) => m.id === id);
    if (mitra) {
      router.push(`/mitra/${id}`)
    }
  };

  const handleDeleteSuccess = () => {
    setIsDeleteOpen(false);
    setSelectedMitra(null);
    router.refresh(); // Refresh data tabel
  };

  const mitraHeaders = [
    { key: 'name', label: 'Name' },
    { key: 'phone_number', label: 'Phone Number' },
    { key: 'outlets', label: 'Outlet Count' },
    { key: 'joinDate', label: 'Join Date' },
  ];

  const mitras = mitraData.map(m => ({
    id: m.id,
    name: m.name,
    phone: m.phone_number, 
    outlets: `${m.outlet_count} ${m.outlet_count !== 1 ? "Outlets" : "Outlet"}`,
    joinDate: new Date(m.created_at).toLocaleDateString("en-US", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      timeZone: "Asia/Jakarta",
    }),
  }));

  return (
    <div className={styles.pageContainer}>
      <div className={styles.headerSection}>
        <div className={styles.titleGroup}>
          <h1>Mitra List</h1>
          <p>Manage and monitor all your partner barbershops.</p>
        </div>
        <button className={styles.btnCreate} onClick={() => setShowModal(true)}>
          <Plus size={18} />
          <span>Create Mitra</span>
        </button>
      </div>

      <div className={styles.statCard}>
        <StatCard 
          label="Total Mitra" 
          value={summaryData.total_mitras.toString()} 
          icon={<Wallet size={12}/>} 
        />
        <StatCard 
          label="Total Cabang" 
          value={summaryData.total_outlets.toString()} 
          icon={<Percent size={12}/>} 
        />
        <StatCard 
          label="Revenue Mitra" 
          value={summaryData.today_revenue_formatted} 
          icon={<ShoppingCart size={12}/>} 
        />
      </div>

      <Table 
        headers={mitraHeaders} 
        data={mitras}
        onDetail={handleOpenView}
        onUpdate={handleOpenUpdate}
        onDelete={handleOpenDelete}
      />

      <Modal 
        isOpen={showModal} 
        onClose={() => setShowModal(false)} 
        title="Register New Mitra"
      >
        <MitraForm onSubmitSuccess={handleRefresh} />
      </Modal>

      <Modal 
        isOpen={isDeleteOpen} 
        onClose={() => setIsDeleteOpen(false)} 
        title="Konfirmasi Hapus"
      >
        {selectedMitra && (
          <DeleteConfirm 
            id={selectedMitra.id}
            name={selectedMitra.name}
            onCancel={() => setIsDeleteOpen(false)}
            onSuccess={handleDeleteSuccess}
          />
        )}
      </Modal>

      <Modal 
        isOpen={isUpdateOpen} 
        onClose={() => setIsUpdateOpen(false)} 
        title="Update Mitra Information"
      >
        {editData && (
          <MitraUpdateForm
            initialData={editData} 
            onSubmitSuccess={handleUpdateSuccess} 
          />
        )}
      </Modal>
    </div>
  );
}
