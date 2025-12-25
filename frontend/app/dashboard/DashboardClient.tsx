"use client";
import React, { useState } from 'react';
import { 
  Wallet, TrendingUp, ShoppingCart, Users, 
  MapPin, Percent, Package, Scissors, Calendar 
} from 'lucide-react';
import { StatCard } from '@/components/statcard/StatCard';
import { LineChart } from '@/components/linechart/LineChart';
import { PieChartComp } from '@/components/piechart/PieChart';
import styles from './Dashboard.module.css';

export default function DashboardClient() {
  // Inisialisasi state tanggal
  const [startDate, setStartDate] = useState(new Date().toISOString().split('T')[0]);
  const [endDate, setEndDate] = useState(new Date().toISOString().split('T')[0]);

  return (
    <div className={styles.pageWrapper}>
      <header className={styles.header}>
        <h1>Dashboard Overview</h1>
        
        {/* Tambahan Fitur Range Tanggal */}
        <div className={styles.filterGroup}>
          <div className={styles.filterIcon}>
            <Calendar size={16} />
          </div>
          <div className={styles.dateInputs}>
            <input 
              type="date" 
              className={styles.dateInput} 
              value={startDate} 
              onChange={(e) => setStartDate(e.target.value)} 
            />
            <span className={styles.dateDivider}>to</span>
            <input 
              type="date" 
              className={styles.dateInput} 
              value={endDate} 
              onChange={(e) => setEndDate(e.target.value)} 
            />
          </div>
        </div>
      </header>

      <section className={styles.statsGrid}>
        <StatCard label="Partner Revenue" value="Rp 12.5M" icon={<Wallet size={18}/>} />
        <StatCard label="Profit (20%)" value="Rp 2.5M" icon={<Percent size={18}/>} />
        <StatCard label="Product Sales" value="Rp 4.2M" icon={<ShoppingCart size={18}/>} />
        <StatCard label="Net Profit" value="Rp 1.1M" icon={<TrendingUp size={18}/>} />
        <StatCard label="Total Mitra" value="48" icon={<Users size={18}/>} />
        <StatCard label="Total Outlet" value="152" icon={<MapPin size={18}/>} />
      </section>

      <main className={styles.mainGrid}>
        <LineChart
          title="Daily Revenue Flow" 
          data={[
            {day: 'Mon', p: 3000, m: 2000}, 
            {day: 'Tue', p: 4000, m: 2500},
            {day: 'Wed', p: 3500, m: 3800},
            {day: 'Thu', p: 5000, m: 3000},
          ]} 
          lines={[
            {key: 'p', color: '#0f172a', name: 'Partner'}, 
            {key: 'm', color: '#38bdf8', name: 'Product'}
          ]} 
        />
        <div className={styles.sideCharts}>
          <PieChartComp title="Products" data={[{name: 'Oil', value: 40}, {name: 'Cream', value:70}]} icon={<Package size={16}/>} />
          <PieChartComp title="Services" data={[{name: 'Cut', value: 60}, {name: 'Wash', value:70}]} icon={<Scissors size={16}/>} />
        </div>
      </main>
    </div>
  );
}