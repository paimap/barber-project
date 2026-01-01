'use client'

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';
import Table from '@/components/table/Table';
import Modal from '@/components/modal/Modal';
import ProductForm from '@/components/forms/product/ProductForm';
import styles from './Services.module.css';
import { ServicesClientProps } from './types';
import DeleteConfirm from '@/components/forms/product/DeleteConfirm';
import ProductUpdateForm from '@/components/forms/product/ProductUpdateForm';

export default function ServicesClient({ productData }: ServicesClientProps) {
  const [showModal, setShowModal] = useState(false);
  const [isMounted, setIsMounted] = useState(false);
  const router = useRouter();
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [selectedProduct, setSelectedMitra] = useState<{id: any, name: string} | null>(null);
  const [isUpdateOpen, setIsUpdateOpen] = useState(false);
  const [editData, setEditData] = useState<any>(null);

  const handleOpenUpdate = (id: any) => {
    const target = productData.find((m: any) => m.ID === id);
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

  useEffect(() => {
    setIsMounted(true);
  }, []);

  const handleRefresh = () => {
    setShowModal(false);
    router.refresh(); 
  };

  const handleOpenDelete = (id: any) => {
    const product = productData.find((m) => m.ID === id);
    if (product) {
      setSelectedMitra({ id: product.ID, name: product.Name });
      setIsDeleteOpen(true);
    }
  };

  const handleDeleteSuccess = () => {
    setIsDeleteOpen(false);
    setSelectedMitra(null);
    router.refresh(); // Refresh data tabel
  };

  const commonHeaders = [
    { key: 'name', label: 'Item Name' },
    { key: 'price', label: 'Price' },
  ];

  // Data statis untuk Services
  const services = [
    { id: 1, name: "Haircut Premium", price: "Rp 50.000" },
    { id: 2, name: "Shaving & Massage", price: "Rp 30.000" },
  ];

  const products = productData.map(m => ({
    id: m.ID,
    name: m.Name,
    price: isMounted 
      ? new Intl.NumberFormat('id-ID', {
          style: 'currency',
          currency: 'IDR',
          minimumFractionDigits: 0,
        }).format(m.Price)
      : 'Rp ...', 
  }));
  
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
          <button className={styles.btnAddSmall} onClick={() => alert('Fitur tambah service segera hadir')}>
            <Plus size={16} /> Add Service
          </button>
        </div>
        <Table 
          headers={commonHeaders} 
          data={services}
          onUpdate={(id) => console.log('Update Service', id)}
          onDelete={handleOpenDelete}
        />
      </section>

      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Shop Products</h2>
          <button className={styles.btnAddSmall} onClick={() => setShowModal(true)}>
            <Plus size={16} /> Add Product
          </button>
        </div>
        <Table 
          headers={commonHeaders} 
          data={products}
          onUpdate={handleOpenUpdate}
          onDelete={handleOpenDelete}
        />
      </section>

      <Modal 
        isOpen={showModal} 
        onClose={() => setShowModal(false)} 
        title="Register New Product"
      >
        <ProductForm onSubmitSuccess={handleRefresh} />
      </Modal>

      <Modal 
        isOpen={isDeleteOpen} 
        onClose={() => setIsDeleteOpen(false)} 
        title="Konfirmasi Hapus"
      >
        {selectedProduct && (
          <DeleteConfirm 
            id={selectedProduct.id}
            name={selectedProduct.name}
            onCancel={() => setIsDeleteOpen(false)}
            onSuccess={handleDeleteSuccess}
          />
        )}
      </Modal>

      <Modal 
        isOpen={isUpdateOpen} 
        onClose={() => setIsUpdateOpen(false)} 
        title="Update Product Information"
      >
        {editData && (
          <ProductUpdateForm
            initialData={editData} 
            onSubmitSuccess={handleUpdateSuccess} 
          />
        )}
      </Modal>
    </div>
  );
}