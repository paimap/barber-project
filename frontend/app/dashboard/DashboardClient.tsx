"use client";

import React from "react";
import {
  Wallet,
  TrendingUp,
  ShoppingCart,
  Percent,
  Package,
  Scissors,
} from "lucide-react";

import { StatCard } from "@/components/statcard/StatCard";
import { LineChart } from "@/components/linechart/LineChart";
import { PieChartComp } from "@/components/piechart/PieChart";
import { Leaderboard } from "@/components/leaderboardcard/Leaderboard";

import styles from "./Dashboard.module.css";

interface DashboardProps {
  initialData: {
    summary: {
      total_revenue: number;
      profit_20: number;
      service_revenue: number;
      product_revenue: number;
    };
    chart: Array<{ day: string; p: number; m: number }>;
    distribution: {
      products: Array<{ name: string; value: number }>;
      services: Array<{ name: string; value: number }>;
    };
    leaderboard: {
      top_outlets: Array<{ name: string; revenue: number }>;
      top_mitras?: Array<{ name: string; revenue: number }>;
      top_barbers?: Array<{ name: string; revenue: number }>; // Data Barber
    };
  };
}

export default function DashboardClient({ initialData }: DashboardProps) {
  const { summary, chart, distribution, leaderboard } = initialData;

  // 🔒 FORMAT AMAN ANTI-HYDRATION
  const formatIDR = (val: number) => {
    return "Rp " + (val ?? 0).toLocaleString("id-ID");
  };

  // Mapping data Cabang
  const branchData = (leaderboard.top_outlets ?? []).map((item) => ({
    name: item.name,
    value: formatIDR(item.revenue),
    subtitle: "Cabang",
  }));

  // Logika Penentuan Kolom Kedua Leaderboard (Barber vs Mitra)
  const hasBarberData = leaderboard.top_barbers && leaderboard.top_barbers.length > 0;
  
  const secondaryLeaderboardData = hasBarberData 
    ? (leaderboard.top_barbers ?? []).map((item) => ({
        name: item.name,
        value: formatIDR(item.revenue),
        subtitle: "Barber",
      }))
    : (leaderboard.top_mitras ?? []).map((item) => ({
        name: item.name,
        value: formatIDR(item.revenue),
        subtitle: "Mitra",
      }));

  const secondaryTitle = hasBarberData ? "Barber Terbaik" : "Mitra Terbaik";
  const secondaryType = hasBarberData ? "barber" : "partner";

  return (
    <div className={styles.pageWrapper}>
      {/* HEADER */}
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h1>Ringkasan Dashboard</h1>
          <p>Ringkasan performa 7 hari terakhir</p>
        </div>
      </header>

      {/* STAT CARDS */}
      <div className={styles.statCard}>
        <StatCard
          label="Total Pendapatan"
          value={formatIDR(summary.total_revenue)}
          icon={<Wallet size={14} />}
        />
        <StatCard
          label="Laba (20%)"
          value={formatIDR(summary.profit_20)}
          icon={<Percent size={14} />}
        />
        <StatCard
          label="Pendapatan Jasa"
          value={formatIDR(summary.service_revenue)}
          icon={<ShoppingCart size={14} />}
        />
        <StatCard
          label="Pendapatan Produk"
          value={formatIDR(summary.product_revenue)}
          icon={<TrendingUp size={14} />}
        />
      </div>

      {/* CHART SECTION */}
      <div className={styles.linePie}>
        <div className={styles.line}>
          {chart.length > 0 ? (
            <LineChart
              title="Aliran Pendapatan Harian"
              data={chart}
              lines={[
                {key: 'm', color: '#0f172a', name: 'Pendapatan Jasa'}, 
                {key: 'p', color: '#38bdf8', name: 'Pendapatan Produk'}
              ]} 
            />
          ) : (
            <div className={styles.emptyState}>
              Belum ada data grafik
            </div>
          )}
        </div>

        <div className={styles.pie}>
          {distribution.products.length > 0 ? (
            <PieChartComp
              title="Distribusi Produk"
              data={distribution.products}
              icon={<Package size={16} />}
            />
          ) : (
            <div className={styles.emptyState}>Data produk kosong</div>
          )}

          {distribution.services.length > 0 ? (
            <PieChartComp
              title="Distribusi Jasa"
              data={distribution.services}
              icon={<Scissors size={16} />}
            />
          ) : (
            <div className={styles.emptyState}>Data jasa kosong</div>
          )}
        </div>
      </div>

      {/* LEADERBOARD */}
      <div className={styles.leaderboardGrid}>
        <Leaderboard
          title="Cabang Terlaris"
          type="branch"
          data={branchData}
        />
        {/* Kolom ini akan otomatis berubah sesuai ketersediaan data top_barbers */}
        <Leaderboard
          title={secondaryTitle}
          type={secondaryType}
          data={secondaryLeaderboardData}
        />
      </div>
    </div>
  );
}