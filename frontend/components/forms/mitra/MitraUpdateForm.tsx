'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface MitraUpdateFormProps {
  initialData: {
    id: any;
    name: string;
    phone_number: string;
  };
  onSubmitSuccess: () => void;
}

export default function MitraUpdateForm({ initialData, onSubmitSuccess }: MitraUpdateFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const data = Object.fromEntries(formData.entries());

    try {
      const res = await fetch(`/api/mitra/${initialData.id}`, {
        method: 'PUT', // Menggunakan PUT atau PATCH
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        onSubmitSuccess();
      } else {
        alert("Failed to update mitra");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.inputGroup}>
        <label>Full Name</label>
        <input name="name" type="text" defaultValue={initialData.name} required />
      </div>
      <div className={styles.inputGroup}>
        <label>Phone Number</label>
        <input name="phone" type="text" defaultValue={initialData.phone_number} required />
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Updating...' : 'Save Changes'}
        </button>
      </div>
    </form>
  );
}