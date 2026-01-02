'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface Barber {
  ID: number;
  Name: string;
  PhoneNumber: string;
}

interface BarberUpdateFormProps {
  initialData: Barber;
  onSubmitSuccess: () => void;
}

export default function BarberUpdateForm({ initialData, onSubmitSuccess }: BarberUpdateFormProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [name, setName] = useState(initialData.Name);
  const [phone, setPhone] = useState(initialData.PhoneNumber);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch(`/api/barber/${initialData.ID}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, phone }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to update barber');
      }

      onSubmitSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      {error && <div className={styles.errorMessage}>{error}</div>}
      
      <div className={styles.inputGroup}>
        <label>Nama</label>
        <input
          type="text"
          placeholder="John Doe"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          disabled={loading}
        />
      </div>
      
      <div className={styles.inputGroup}>
        <label>Nomor Telepon</label>
        <input
          type="tel"
          placeholder="+62812345678"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          required
          disabled={loading}
        />
      </div>

      <div className={styles.footer}>
        <button
          type="submit"
          disabled={loading}
          className={styles.btnPrimary}
        >
          {loading ? 'Updating...' : 'Update Barber'}
        </button>
      </div>
    </form>
  );
}
