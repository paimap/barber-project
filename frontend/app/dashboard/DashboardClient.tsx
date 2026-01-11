"use client";
import React, { useState } from 'react';
import { 
  Wallet, TrendingUp, ShoppingCart, Users, 
  MapPin, Percent, Package, Scissors, Calendar,
  Trophy, Medal
} from 'lucide-react';
import { StatCard } from '@/components/statcard/StatCard';
import { LineChart } from '@/components/linechart/LineChart';
import { PieChartComp } from '@/components/piechart/PieChart';
import styles from './Dashboard.module.css';
import { Leaderboard } from '@/components/leaderboardcard/Leaderboard';

export default function DashboardClient() {
  const branchData = [
    { name: "Jakarta Pusat", value: "Rp 45.0M", growth: "+12%", subtitle: "Main Branch" },
    { name: "Bandung Core", value: "Rp 38.5M", growth: "+8%", subtitle: "Sarijadi" },
    { name: "Surabaya West", value: "Rp 32.2M", growth: "+5%", subtitle: "Pakuwon" },
  ];

  const partnerData = [
    { name: "Budi Santoso", value: "150 Sales", subtitle: "Level: Gold" },
    { name: "Siti Aminah", value: "142 Sales", subtitle: "Level: Silver" },
    { name: "Andi Wijaya", value: "128 Sales", subtitle: "Level: Silver" },
  ];
  return (
    <div className={styles.pageWrapper}>
      {/* Header dengan kata-kata di bawahnya */}
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h1>Dashboard Overview</h1>
          <p>Selamat datang kembali! Berikut adalah ringkasan performa hari ini.</p>
        </div>
      </header>

      <div className={styles.statCard}>
        <StatCard label="Revenue" value="Rp 12.500.000" icon={<Wallet size={12}/>} />
        <StatCard label="Profit (20%)" value="Rp 2.500.000" icon={<Percent size={12}/>} />
        <StatCard label="Product Sales" value="15" icon={<ShoppingCart size={12}/>} />
        <StatCard label="Total Service" value="25" icon={<TrendingUp size={12}/>} />
      </div>
      <div className={styles.linePie}>
        <div className={styles.line}>
          <LineChart
            title="Daily Revenue Flow" 
            data={[
              {day: 'Mon', p: 3000, m: 2000}, 
              {day: 'Tue', p: 4000, m: 2500},
              {day: 'Wed', p: 3500, m: 3800},
              {day: 'Thu', p: 5000, m: 3000},
              {day: 'Fri', p: 5000, m: 3000},
              {day: 'Sat', p: 5000, m: 7000},
              {day: 'Sun', p: 5000, m: 3000},
            ]} 
            lines={[
              {key: 'p', color: '#0f172a', name: 'Partner'}, 
              {key: 'm', color: '#38bdf8', name: 'Product'}
            ]} 
          />
        </div>
        <div className={styles.pie}>
            <PieChartComp title="Products" data={[{name: 'Oil', value: 40}, {name: 'Cream', value:70}]} icon={<Package size={16}/>} />
            <PieChartComp title="Services" data={[{name: 'Cut', value: 60}, {name: 'Wash', value:70}]} icon={<Scissors size={16}/>} />
        </div>
      </div>

      {/* SECTION BARU: Leaderboard */}
      <div className={styles.leaderboardGrid}>
        <Leaderboard 
          title="Cabang Terlaris" 
          type="branch" 
          data={branchData} 
        />
        <Leaderboard 
          title="Mitra Terbaik" 
          type="partner" 
          data={partnerData} 
        />
      </div>
    </div>
  );
}