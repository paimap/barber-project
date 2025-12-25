import React from 'react';
import styles from './LeaderboardCard.module.css';
import { Trophy } from 'lucide-react';

interface LeaderboardItem {
  name: string;
  subName?: string;
  total: string;
}

interface LeaderboardCardProps {
  title: string;
  icon: React.ReactNode;
  data: LeaderboardItem[];
}

export const LeaderboardCard = ({ title, icon, data }: LeaderboardCardProps) => (
  <div className={styles.card}>
    <div className={styles.header}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        {icon}
        <h2>{title}</h2>
      </div>
      <Trophy size={18} color="#f59e0b" />
    </div>
    <div className={styles.list}>
      {data.map((item, index) => (
        <div key={index} className={styles.item}>
          <div className={styles.itemInfo}>
            <div className={`${styles.rank} ${
              index === 0 ? styles.top1 : index === 1 ? styles.top2 : index === 2 ? styles.top3 : styles.otherRank
            }`}>
              {index + 1}
            </div>
            <div className={styles.nameContainer}>
              <span className={styles.name}>{item.name}</span>
              {item.subName && <span className={styles.subName}>{item.subName}</span>}
            </div>
          </div>
          <span className={styles.value}>{item.total}</span>
        </div>
      ))}
    </div>
  </div>
);