"use client";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts';
import styles from './LineChart.module.css';

interface RevenueChartProps {
  title: string;
  data: any[];
  lines: { key: string; color: string; name: string }[];
}

export const LineChart = ({ title, data, lines }: RevenueChartProps) => {
  
  // Helper untuk format angka ke IDR yang ringkas (Contoh: 1.2jt)
  const formatYAxis = (value: number) => {
    if (value >= 1000000) return `${(value / 1000000).toFixed(1)}jt`;
    if (value >= 1000) return `${(value / 1000).toFixed(0)}rb`;
    return value.toString();
  };

  // Helper untuk format IDR lengkap di Tooltip
  const formatTooltip = (value: number) => 
    new Intl.NumberFormat('id-ID', { 
      style: 'currency', 
      currency: 'IDR', 
      maximumFractionDigits: 0 
    }).format(value);

  return (
    <div className={styles.chartContainer}>
      <h2 className={styles.title}>{title}</h2>
      <ResponsiveContainer width="100%" height={400}>
        <AreaChart data={data}>
          <defs>
            {lines.map((line) => (
              <linearGradient key={line.key} id={`grad${line.key}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor={line.color} stopOpacity={0.1}/>
                <stop offset="95%" stopColor={line.color} stopOpacity={0}/>
              </linearGradient>
            ))}
          </defs>
          <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
          <XAxis 
            dataKey="day" 
            axisLine={false} 
            tickLine={false} 
            tick={{fill: '#94a3b8', fontSize: 11}} 
          />
          <YAxis 
            axisLine={false} 
            tickLine={false} 
            tick={{fill: '#94a3b8', fontSize: 11}}
            tickFormatter={formatYAxis} // Menambahkan format di sumbu Y
            width={60}
          />
          <Tooltip 
            formatter={(value: any) => {
              // Pastikan kita selalu mengembalikan string atau number, bukan undefined
              const val = Number(value);
              if (isNaN(val)) return ["Rp 0", ""];
              
              return [formatTooltip(val), ""];
            }}
            contentStyle={{ 
              borderRadius: '12px', 
              border: 'none', 
              boxShadow: '0 10px 30px rgba(0,0,0,0.1)' 
            }} 
          />
          <Legend 
            verticalAlign="top" 
            align="right" 
            iconType="circle" 
            wrapperStyle={{ fontSize: '12px', paddingBottom: '20px' }} 
          />
          {lines.map((line) => (
            <Area 
              key={line.key} 
              type="monotone" 
              dataKey={line.key} 
              stroke={line.color} 
              strokeWidth={3} 
              fill={`url(#grad${line.key})`} 
              name={line.name} 
            />
          ))}
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};