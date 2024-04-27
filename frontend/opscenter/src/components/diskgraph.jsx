import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import moment from 'moment';

export const DiskGraph = ({ data }) => {
    return (
      <ResponsiveContainer width="100%" height={300}>
        <LineChart
          data={data}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="tps" stroke="#8884d8" name="TPS" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="KBPerTrans" stroke="#82ca9d" name="KB/Trans" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="MBPerSec" stroke="#ffc658" name="MB/s" activeDot={{ r: 8 }} />
        </LineChart>
      </ResponsiveContainer>
    );
  };
  