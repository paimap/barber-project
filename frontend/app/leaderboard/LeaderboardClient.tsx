"use client";
import React, { useState } from 'react';
import { Package, Scissors, Users, MapPin, Calendar } from 'lucide-react';
import { LeaderboardCard } from '@/components/leaderboardcard/LeaderboardCard';
import styles from './Leaderboard.module.css';

export default function LeaderboardClient() {
  const [startDate, setStartDate] = useState("2024-01-01");
  const [endDate, setEndDate] = useState("2024-01-31");

  // Mock Data
  const topProducts = [
    { name: "Matte Pomade", total: "1,240 Sold" },
    { name: "Clay Wax", total: "850 Sold" },
    { name: "Beard Oil Premium", total: "420 Sold" },
  ];

  const topServices = [
    { name: "Gentleman Haircut", total: "2,100 Svcs" },
    { name: "Hot Towel Shave", total: "1,150 Svcs" },
    { name: "Hair Coloring", total: "340 Svcs" },
  ];

  const topPartners = [
    { name: "Budi Santoso", subName: "Jakarta Selatan", total: "Rp 45.2M" },
    { name: "Andi Wijaya", subName: "Bandung", total: "Rp 38.7M" },
    { name: "Siti Aminah", subName: "Surabaya", total: "Rp 32.1M" },
  ];

  const topBranches = [
    { name: "Barberia Tebet", subName: "Owner: Budi Santoso", total: "Rp 28.4M" },
    { name: "Barberia Dago", subName: "Owner: Andi Wijaya", total: "Rp 22.1M" },
    { name: "Barberia Kemang", subName: "Owner: Budi Santoso", total: "Rp 19.8M" },
  ];

  return (
    <div className={styles.pageWrapper}>
      <header className={styles.header}>
        <h1>Ranking Leaderboard</h1>
        <div className={styles.filterGroup}>
          <Calendar size={16} color="#64748b" />
          <input type="date" className={styles.dateInput} value={startDate} onChange={(e) => setStartDate(e.target.value)} />
          <span>to</span>
          <input type="date" className={styles.dateInput} value={endDate} onChange={(e) => setEndDate(e.target.value)} />
        </div>
      </header>

      <main className={styles.boardGrid}>
        <LeaderboardCard title="Top Products" icon={<Package size={20} color="#0f172a"/>} data={topProducts} />
        <LeaderboardCard title="Top Services" icon={<Scissors size={20} color="#0f172a"/>} data={topServices} />
        <LeaderboardCard title="Top Partners" icon={<Users size={20} color="#0f172a"/>} data={topPartners} />
        <LeaderboardCard title="Top Branches" icon={<MapPin size={20} color="#0f172a"/>} data={topBranches} />
      </main>
    </div>
  );
}