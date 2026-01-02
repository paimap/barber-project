'use client'

import { useState } from 'react';
import styles from '@/components/forms/Form.module.css';

interface BarberFormProps {
  onSubmitSuccess: () => void;
}

export default function BarberForm({ onSubmitSuccess }: BarberFormProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/barber', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, name, phone }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to create barber');
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
        <label>Email</label>
        <input
          type="email"
          placeholder="barber@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          disabled={loading}
        />
      </div>
      
      <div className={styles.inputGroup}>
        <label>Password</label>
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          disabled={loading}
        />
      </div>
      
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
          {loading ? 'Creating...' : 'Create Barber'}
        </button>
      </div>
    </form>
  );
}
