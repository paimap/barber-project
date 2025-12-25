"use client";
import React from 'react';
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import styles from './PieChart.module.css';

interface DistributionCardProps {
  title: string;
  data: any[];
  icon: React.ReactNode;
}

const COLORS = ['#0f172a', '#38bdf8', '#94a3b8'];

export const PieChartComp = ({ title, data, icon }: DistributionCardProps) => (
  <div className={styles.card}>
    <div className={styles.header}>
      {icon}
      <h2>{title}</h2>
    </div>
    <ResponsiveContainer width="100%" height={150}>
      <PieChart>
        <Pie data={data} innerRadius={45} outerRadius={60} paddingAngle={5} dataKey="value" stroke="none">
          {data.map((_, i) => <Cell key={i} fill={COLORS[i % COLORS.length]} />)}
        </Pie>
        <Tooltip contentStyle={{ borderRadius: '8px', border: 'none' }} />
        <Legend verticalAlign="middle" align="right" layout="vertical" iconSize={8} wrapperStyle={{ fontSize: '11px' }} />
      </PieChart>
    </ResponsiveContainer>
  </div>
);