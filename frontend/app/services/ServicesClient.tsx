'use client'

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';
import Table from '@/components/table/Table';
import Modal from '@/components/modal/Modal';
import styles from './Services.module.css';
import { ServicesClientProps } from './types';

// Components
import ProductForm from '@/components/forms/product/ProductForm';
import ProductUpdateForm from '@/components/forms/product/ProductUpdateForm';
import DeleteConfirm from '@/components/forms/product/DeleteConfirm';

import ServiceTypeForm from '@/components/forms/service-type/ServiceTypeForm';
import ServiceTypeUpdateForm from '@/components/forms/service-type/ServiceTypeUpdateForm';
import DeleteConfirmServiceType from '@/components/forms/service-type/DeleteConfirmServiceType';

export default function ServicesClient({ productData, serviceData }: ServicesClientProps) {
  const [isMounted, setIsMounted] = useState(false);
  const router = useRouter();

  // --- STATE UNTUK PRODUCT ---
  const [showCreateProduct, setShowCreateProduct] = useState(false);
  const [showUpdateProduct, setShowUpdateProduct] = useState(false);
  const [showDeleteProduct, setShowDeleteProduct] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<any>(null);

  // --- STATE UNTUK SERVICE ---
  const [showCreateService, setShowCreateService] = useState(false);
  const [showUpdateService, setShowUpdateService] = useState(false);
  const [showDeleteService, setShowDeleteService] = useState(false);
  const [selectedService, setSelectedService] = useState<any>(null);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  // --- HANDLERS PRODUCT ---
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

  // --- HANDLERS SERVICE ---
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
    { key: 'name', label: 'Item Name' },
    { key: 'price', label: 'Price' },
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
          <h1>Services & Products</h1>
          <p>Configure your barber offerings and shop inventory.</p>
        </div>
      </header>

      {/* SECTION: SERVICES */}
      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Barber Services</h2>
          <button className={styles.btnAddSmall} onClick={() => setShowCreateService(true)}>
            <Plus size={16} /> Add Service
          </button>
        </div>
        <Table 
          headers={commonHeaders} 
          data={formatData(serviceData)}
          onUpdate={handleOpenUpdateService}
          onDelete={handleOpenDeleteService}
        />
      </section>

      {/* SECTION: PRODUCTS */}
      <section className={styles.section}>
        <div className={styles.sectionTitle}>
          <h2>Products</h2>
          <button className={styles.btnAddSmall} onClick={() => setShowCreateProduct(true)}>
            <Plus size={16} /> Add Product
          </button>
        </div>
        <Table 
          headers={commonHeaders} 
          data={formatData(productData)}
          onUpdate={handleOpenUpdateProduct}
          onDelete={handleOpenDeleteProduct}
        />
      </section>

      {/* --- MODALS PRODUCT --- */}
      <Modal isOpen={showCreateProduct} onClose={() => setShowCreateProduct(false)} title="Register New Product">
        <ProductForm onSubmitSuccess={onProductSuccess} />
      </Modal>

      <Modal isOpen={showUpdateProduct} onClose={() => setShowUpdateProduct(false)} title="Update Product">
        {selectedProduct && <ProductUpdateForm initialData={selectedProduct} onSubmitSuccess={onProductSuccess} />}
      </Modal>

      <Modal isOpen={showDeleteProduct} onClose={() => setShowDeleteProduct(false)} title="Hapus Produk">
        {selectedProduct && (
          <DeleteConfirm id={selectedProduct.id} name={selectedProduct.name} onCancel={() => setShowDeleteProduct(false)} onSuccess={onProductSuccess} />
        )}
      </Modal>

      {/* --- MODALS SERVICE --- */}
      <Modal isOpen={showCreateService} onClose={() => setShowCreateService(false)} title="Add New Service">
        <ServiceTypeForm onSubmitSuccess={onServiceSuccess} />
      </Modal>

      <Modal isOpen={showUpdateService} onClose={() => setShowUpdateService(false)} title="Update Service">
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