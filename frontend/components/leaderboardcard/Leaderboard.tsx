import React from 'react';
import { MapPin, Trophy, Medal, User, LucideIcon } from 'lucide-react';
import styles from './Leaderboard.module.css';

interface LeaderboardItem {
  name: string;
  value: string | number;
}

interface LeaderboardProps {
  title: string;
  // Tambahkan 'barber' ke dalam union type
  type: 'branch' | 'partner' | 'barber'; 
  data: LeaderboardItem[];
}

export const Leaderboard: React.FC<LeaderboardProps> = ({ title, type, data }) => {
  // Menentukan Icon berdasarkan tipe secara dinamis
  const getHeaderIcon = (): LucideIcon => {
    switch (type) {
      case 'branch': return MapPin;
      case 'barber': return User; // Icon khusus untuk Barber
      default: return Trophy;     // Default untuk Partner/Mitra
    }
  };

  const HeaderIcon = getHeaderIcon();

  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <div className={styles.titleGroup}>
          <HeaderIcon size={18} className={styles.icon} />
          <h3>{title}</h3>
        </div>
      </div>
      
      <div className={styles.list}>
        {data.length > 0 ? (
          data.map((item, i) => (
            <div key={i} className={styles.listItem}>
              <div className={styles.itemLeft}>
                <div className={styles.rankContainer}>
                  {i === 0 ? (
                    <Medal size={20} className={styles.goldMedal} fill="currentColor" />
                  ) : (
                    <span className={styles.rankNumber}>{i + 1}</span>
                  )}
                </div>
                <div className={styles.info}>
                  <p className={styles.name}>{item.name}</p>
                </div>
              </div>

              <div className={styles.itemRight}>
                <span className={styles.mainValue}>{item.value}</span>
              </div>
            </div>
          ))
        ) : (
          <div className={styles.emptyText}>Tidak ada data ranking</div>
        )}
      </div>
    </div>
  );
};