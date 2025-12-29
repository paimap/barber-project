'use client'

import { useState } from 'react';
import styles from './Form.module.css';

interface MitraFormProps {
  onSubmitSuccess: () => void;
}

export default function MitraForm({ onSubmitSuccess }: MitraFormProps) {
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const data = Object.fromEntries(formData.entries());

    try {
      const res = await fetch('/api/mitra', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        onSubmitSuccess();
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
        <input name="name" type="text" placeholder="e.g. John Doe" required />
      </div>
      <div className={styles.inputGroup}>
        <label>Email Address</label>
        <input name="email" type="email" placeholder="john@example.com" required />
      </div>
      <div className={styles.inputGroup}>
        <label>Password</label>
        <input name="password" type="password" placeholder="••••••••" required />
      </div>
      <div className={styles.inputGroup}>
        <label>Phone Number</label>
        <input name="phone" type="text" placeholder="0812..." required />
      </div>
      
      <div className={styles.footer}>
        <button type="submit" className={styles.btnPrimary} disabled={loading}>
          {loading ? 'Saving...' : 'Create Mitra'}
        </button>
      </div>
    </form>
  );
}