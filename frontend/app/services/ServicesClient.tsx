'use client'

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Plus, Wallet, Percent } from 'lucide-react';
import { useAuth } from '@/context/AuthContext'; 

import Table from '@/components/table/Table';
import Modal from '@/components/modal/Modal';
import styles from './Services.module.css';
import { ServicesClientProps } from './types';
import { StatCard } from '@/components/statcard/StatCard';

// Komponen Form
import ProductForm from '@/components/forms/product/ProductForm';
import ProductUpdateForm from '@/components/forms/product/ProductUpdateForm';
import DeleteConfirm from '@/components/forms/product/DeleteConfirm';

import ServiceTypeForm from '@/components/forms/service-type/ServiceTypeForm';
import ServiceTypeUpdateForm from '@/components/forms/service-type/ServiceTypeUpdateForm';
import DeleteConfirmServiceType from '@/components/forms/service-type/DeleteConfirmServiceType';

export default function ServicesClient({ productData, serviceData, summaryData }: ServicesClientProps) {
  const { user } = useAuth(); 
  const [isMounted, setIsMounted] = useState(false);
  const router = useRouter();

  // --- STATE UNTUK PRODUK ---
  const [showCreateProduct, setShowCreateProduct] = useState(false);
  const [showUpdateProduct, setShowUpdateProduct] = useState(false);
  const [showDeleteProduct, setShowDeleteProduct] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<any>(null);

  // --- STATE UNTUK LAYANAN ---
  const [showCreateService, setShowCreateService] = useState(false);
  const [showUpdateService, setShowUpdateService] = useState(false);
  const [showDeleteService, setShowDeleteService] = useState(false);
  const [selectedService, setSelectedService] = useState<any>(null);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Cek apakah user adalah Superadmin
  const isSuperAdmin = user?.role === "SUPERADMIN";

  const formatIDR = (val: number) => 
    new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(val);

  // --- HANDLERS PRODUK ---
  const handleOpenUpdateProduct = (id: any) => {
    const target = productData.find((m: any) => m.ID === id);
    if (target) {
      setSelectedProduct({ id: target.ID, name: target.Name, price: target.Price });
      setShowUpdateProduct(true);
    }
  };

  const handleOpenDeleteProduct = (id: any) => {
    const target = productData.find((m) => m.ID === id);
    if (target) {
      setSelectedProduct({ id: target.ID, name: target.Name });
      setShowDeleteProduct(true);
    }
  };

  const onProductSuccess = () => {
    setShowCreateProduct(false);
    setShowUpdateProduct(false);
    setShowDeleteProduct(false);
    setSelectedProduct(null);
    router.refresh();
  };

  // --- HANDLERS LAYANAN ---
  const handleOpenUpdateService = (id: any) => {
    const target = serviceData.find((m: any) => m.ID === id);
    if (target) {
      setSelectedService({ id: target.ID, name: target.Name, price: target.Price });
      setShowUpdateService(true);
    }
  };

  const handleOpenDeleteService = (id: any) => {
    const target = serviceData.find((m) => m.ID === id);
    if (target) {
      setSelectedService({ id: target.ID, name: target.Name });
      setShowDeleteService(true);
    }
  };

  const onServiceSuccess = () => {
    setShowCreateService(false);
    setShowUpdateService(false);
    setShowDeleteService(false);
    setSelectedService(null);
    router.refresh();
  };

  // --- FORMATTING ---
  const commonHeaders = [
    { key: 'name', label: 'Nama Item' },
    { key: 'price', label: 'Harga' },
  ];

  const formatData = (data: any[]) => data.map(m => ({
    id: m.ID,
    name: m.Name,
    price: isMounted 
      ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(m.Price)
      : 'Rp ...', 
  }));

  return (
    <div className={styles.pageContainer}>
      <header className={styles.headerSection}>
        <div className={styles.titleGroup}>
          <h1>Layanan & Produk</h1>
          <p>Konfigurasi penawaran barber dan inventaris toko Anda.</p>
        </div>
      </header>

      <div className={styles.statCard}>
        <StatCard 
            label="Pendapatan Produk" 
            value={isMounted ? formatIDR(summaryData.product_revenue) : "Rp ..."} 
            icon={<Wallet size={12}/>} 
        />
        <StatCard 
            label="Pendapatan Jasa" 
            value={isMounted ? formatIDR(summaryData.service_revenue) : "Rp ..."} 
            icon={<Percent size={12}/>} 
        />
        <StatCard 
            label="Total Produk Terjual" 
            value={isMounted ? summaryData.product_sold.toString() : "..."} 
            icon={<Wallet size={12}/>} 
        />
        <StatCard 
            label="Total Jasa Dilakukan" 
            value={isMounted ? summaryData.service_performed.toString() : "..."} 
            icon={<Percent size={12}/>} 
        />
      </div>

      {/* SEKSI: LAYANAN */}
      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Layanan Barber</h2>
          {/* HANYA MUNCUL JIKA SUPERADMIN */}
          {isSuperAdmin && (
            <button className={styles.btnAddSmall} onClick={() => setShowCreateService(true)}>
              <Plus size={16} /> Tambah Layanan
            </button>
          )}
        </div>
        <Table 
          headers={commonHeaders} 
          data={formatData(serviceData)}
          onUpdate={isSuperAdmin ? handleOpenUpdateService : undefined}
          onDelete={isSuperAdmin ? handleOpenDeleteService : undefined}
        />
      </section>

      {/* SEKSI: PRODUK */}
      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Produk</h2>
          {/* HANYA MUNCUL JIKA SUPERADMIN */}
          {isSuperAdmin && (
            <button className={styles.btnAddSmall} onClick={() => setShowCreateProduct(true)}>
              <Plus size={16} /> Tambah Produk
            </button>
          )}
        </div>
        <Table 
          headers={commonHeaders} 
          data={formatData(productData)}
          onUpdate={isSuperAdmin ? handleOpenUpdateProduct : undefined}
          onDelete={isSuperAdmin ? handleOpenDeleteProduct : undefined}
        />
      </section>

      {/* --- MODAL PRODUK --- */}
      <Modal isOpen={showCreateProduct} onClose={() => setShowCreateProduct(false)} title="Daftarkan Produk Baru">
        <ProductForm onSubmitSuccess={onProductSuccess} />
      </Modal>

      <Modal isOpen={showUpdateProduct} onClose={() => setShowUpdateProduct(false)} title="Perbarui Produk">
        {selectedProduct && <ProductUpdateForm initialData={selectedProduct} onSubmitSuccess={onProductSuccess} />}
      </Modal>

      <Modal isOpen={showDeleteProduct} onClose={() => setShowDeleteProduct(false)} title="Hapus Produk">
        {selectedProduct && (
          <DeleteConfirm id={selectedProduct.id} name={selectedProduct.name} onCancel={() => setShowDeleteProduct(false)} onSuccess={onProductSuccess} />
        )}
      </Modal>

      {/* --- MODAL LAYANAN --- */}
      <Modal isOpen={showCreateService} onClose={() => setShowCreateService(false)} title="Tambah Layanan Baru">
        <ServiceTypeForm onSubmitSuccess={onServiceSuccess} />
      </Modal>

      <Modal isOpen={showUpdateService} onClose={() => setShowUpdateService(false)} title="Perbarui Layanan">
        {selectedService && <ServiceTypeUpdateForm initialData={selectedService} onSubmitSuccess={onServiceSuccess} />}
      </Modal>

      <Modal isOpen={showDeleteService} onClose={() => setShowDeleteService(false)} title="Hapus Layanan">
        {selectedService && (
          <DeleteConfirmServiceType id={selectedService.id} name={selectedService.name} onCancel={() => setShowDeleteService(false)} onSuccess={onServiceSuccess} />
        )}
      </Modal>
    </div>
  );
}