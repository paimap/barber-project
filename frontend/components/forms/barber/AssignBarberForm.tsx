'use client'

import { useState, useEffect } from 'react';
import styles from '../Form.module.css';

interface Barber {
  ID: number;
  Name: string;
}

interface Outlet {
  ID: number;
  Address: string;
}

interface AssignBarberFormProps {
  barberData: Barber;
  onSubmitSuccess: () => void;
}

export default function AssignBarberForm({ barberData, onSubmitSuccess }: AssignBarberFormProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [outlets, setOutlets] = useState<Outlet[]>([]);
  const [selectedOutletId, setSelectedOutletId] = useState('');

  useEffect(() => {
    const fetchOutlets = async () => {
      try {
        const response = await fetch('/api/outlets', {
          method: 'GET',
        });
        if (response.ok) {
          const data = await response.json();
          setOutlets(data.outlets || []);
        }
      } catch (err) {
        console.error('Error fetching outlets:', err);
      }
    };

    fetchOutlets();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (!selectedOutletId) {
      setError('Please select an outlet');
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`/api/barber/${barberData.ID}/assign`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ outlet_id: parseInt(selectedOutletId) }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Failed to assign barber');
      }

      onSubmitSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      {error && (
        <div style={{ color: '#dc2626', fontSize: '14px' }}>{error}</div>
      )}

      <div className={styles.inputGroup}>
        <label>Assigning:</label>
        <span style={{ fontWeight: 600 }}>{barberData.Name}</span>
      </div>

      <div className={styles.inputGroup}>
        <label htmlFor="outlet">Outlet</label>
        <select
          id="outlet"
          value={selectedOutletId}
          onChange={(e) => setSelectedOutletId(e.target.value)}
          required
          disabled={loading}
          className={styles.inputGroup}
        >
          <option value="">-- Select Outlet --</option>
          {outlets.map((outlet) => (
            <option key={outlet.ID} value={outlet.ID}>
              {outlet.Address}
            </option>
          ))}
        </select>
      </div>

      <div className={styles.footer}>
        <button
          type="submit"
          disabled={loading}
          className={styles.btnPrimary}
        >
          {loading ? 'Assigning...' : 'Assign to Outlet'}
        </button>
      </div>
    </form>
  );
}
