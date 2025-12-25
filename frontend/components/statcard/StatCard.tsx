import React from 'react';
import styles from './StatCard.module.css';

interface StatCardProps {
  label: string;
  value: string;
  icon: React.ReactNode;
}

export const StatCard = ({ label, value, icon }: StatCardProps) => (
  <div className={styles.card}>
    <div className={styles.iconBox}>{icon}</div>
    <span className={styles.label}>{label}</span>
    <span className={styles.value}>{value}</span>
  </div>
);