import React from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export const CPUGraph = ({ data }) => {
  const formattedData = data.map(item => ({
    time: item.time,
    idle: item.idle,
    user: item.user,
    system: item.system
  }));

  return (
    <ResponsiveContainer width="100%" height={300}>
      <AreaChart
        data={formattedData}
        margin={{
          top: 5,
          right: 30,
          left: 20,
          bottom: 5,
        }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" />
        <YAxis />
        <Tooltip />
        <Area type="monotone" dataKey="idle" stroke="#8884d8" fillOpacity={0.5} fill="#8884d8" />
        <Area type="monotone" dataKey="user" stroke="#82ca9d" fillOpacity={0.5} fill="#82ca9d" />
        <Area type="monotone" dataKey="system" stroke="#ffc658" fillOpacity={0.5} fill="#ffc658" />
      </AreaChart>
    </ResponsiveContainer>
  );
};
