import React from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

export const DiskGraph = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <AreaChart
        data={data}
        margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Area type="monotone" dataKey="tps" stroke="#8884d8" fill="#8884d8" name="TPS" fillOpacity={0.3} />
        <Area type="monotone" dataKey="KBPerTrans" stroke="#82ca9d" fill="#82ca9d" name="KB/Trans" fillOpacity={0.3} />
        <Area type="monotone" dataKey="MBPerSec" stroke="#ffc658" fill="#ffc658" name="MB/s" fillOpacity={0.3} />
      </AreaChart>
    </ResponsiveContainer>
  );
};
