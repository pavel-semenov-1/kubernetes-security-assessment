import React from 'react';
import { StatResult } from '../types/StatResult';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface BarChartProps {
  data: StatResult[];
}

const MisconfigurationBarChart: React.FC<BarChartProps> = ({ data }) => {
  const formattedData = data.map((item) => ({
    tool: item.scanner,
    critical: item.stats.misconfigurations.CRITICAL,
    high: item.stats.misconfigurations.HIGH,
    medium: item.stats.misconfigurations.MEDIUM,
    low: item.stats.misconfigurations.LOW,
  }));

  return (
    <ResponsiveContainer width="100%" height={370}>
      <BarChart data={formattedData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="tool" tick={{ fill: 'white', fontSize: 16 }} />
        <YAxis tick={{ fill: 'white', fontSize: 14 }} />
        <Tooltip />
        <Legend />
        <Bar dataKey="critical" fill="#ff0000" />
        <Bar dataKey="high" fill="#ff8000" />
        <Bar dataKey="medium" fill="#ffff00" />
        <Bar dataKey="low" fill="#00ff00" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default MisconfigurationBarChart;
