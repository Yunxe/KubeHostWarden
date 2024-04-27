import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import moment from 'moment';

export const CPUGraph = ({ data }) => {
  const formatXAxis = (tickItem) => {
    return tickItem; // X轴已经是格式化的时间，无需再次格式化
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart
        data={data}
        margin={{
          top: 5,
          right: 30,
          left: 20,
          bottom: 5,
        }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" tickFormatter={formatXAxis} />
        <YAxis />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="idle" stroke="#8884d8" name="Idle CPU" activeDot={{ r: 8 }} />
        <Line type="monotone" dataKey="user" stroke="#82ca9d" name="User CPU" activeDot={{ r: 8 }} />
        <Line type="monotone" dataKey="system" stroke="#ffc658" name="System CPU" activeDot={{ r: 8 }} />
      </LineChart>
    </ResponsiveContainer>
  );
};
