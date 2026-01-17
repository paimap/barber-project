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
      top_mitras: Array<{ name: string; revenue: number }>;
    };
  };
}

export default function DashboardClient({ initialData }: DashboardProps) {
  const { summary, chart, distribution, leaderboard } = initialData;

  // 🔒 FORMAT AMAN ANTI-HYDRATION
  const formatIDR = (val: number) => {
    return "Rp " + (val ?? 0).toLocaleString("id-ID");
  };

  const branchData = (leaderboard.top_outlets ?? []).map((item) => ({
    name: item.name,
    value: formatIDR(item.revenue),
    subtitle: "Cabang",
  }));

  const partnerData = (leaderboard.top_mitras ?? []).map((item) => ({
    name: item.name,
    value: formatIDR(item.revenue),
    subtitle: "Mitra",
  }));

  return (
    <div className={styles.pageWrapper}>
      {/* HEADER */}
      <header className={styles.header}>
        <div className={styles.headerTitle}>
          <h1>Dashboard Overview</h1>
          <p>Ringkasan performa 7 hari terakhir</p>
        </div>
      </header>

      {/* STAT CARDS */}
      <div className={styles.statCard}>
        <StatCard
          label="Total Revenue"
          value={formatIDR(summary.total_revenue)}
          icon={<Wallet size={14} />}
        />
        <StatCard
          label="Profit (20%)"
          value={formatIDR(summary.profit_20)}
          icon={<Percent size={14} />}
        />
        <StatCard
          label="Service Revenue"
          value={formatIDR(summary.service_revenue)}
          icon={<ShoppingCart size={14} />}
        />
        <StatCard
          label="Product Revenue"
          value={formatIDR(summary.product_revenue)}
          icon={<TrendingUp size={14} />}
        />
      </div>

      {/* CHART SECTION */}
      <div className={styles.linePie}>
        <div className={styles.line}>
          {chart.length > 0 ? (
            <LineChart
              title="Daily Revenue Flow"
              data={chart}
              lines={[
                {key: 'm', color: '#0f172a', name: 'Service Revenue'}, 
                {key: 'p', color: '#38bdf8', name: 'Product revenue'}
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
              title="Products"
              data={distribution.products}
              icon={<Package size={16} />}
            />
          ) : (
            <div className={styles.emptyState}>Data produk kosong</div>
          )}

          {distribution.services.length > 0 ? (
            <PieChartComp
              title="Services"
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
        <Leaderboard
          title="Mitra Terbaik"
          type="partner"
          data={partnerData}
        />
      </div>
    </div>
  );
}
