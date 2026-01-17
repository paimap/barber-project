"use client";
import React from 'react';
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import styles from './PieChart.module.css';

interface DistributionCardProps {
  title: string;
  data: any[];
  icon: React.ReactNode;
}

const COLORS = [
  '#0f172a', // Slate 950
  '#38bdf8', // Sky 400
  '#6366f1', // Indigo 500
  '#a855f7', // Purple 500
  '#10b981', // Emerald 500
  '#f59e0b', // Amber 500
];

export const PieChartComp = ({ title, data, icon }: DistributionCardProps) => (
  <div className={styles.card}>
    <div className={styles.header}>
      {icon}
      <h2>{title}</h2>
    </div>
    <ResponsiveContainer width="100%" height={150}>
      <PieChart>
        <Pie data={data} innerRadius={60} outerRadius={75} paddingAngle={5} dataKey="value" stroke="none">
          {data.map((_, i) => <Cell key={i} fill={COLORS[i % COLORS.length]} />)}
        </Pie>
        <Tooltip contentStyle={{ borderRadius: '8px', border: 'none' }} />
        <Legend verticalAlign="middle" align="right" layout="vertical" iconSize={8} wrapperStyle={{ fontSize: '11px' }} />
      </PieChart>
    </ResponsiveContainer>
  </div>
);