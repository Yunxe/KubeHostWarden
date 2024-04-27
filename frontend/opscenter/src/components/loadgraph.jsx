import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import moment from 'moment';

export const LoadGraph = ({ data }) => {
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
          <Line type="monotone" dataKey="one" stroke="#8884d8" name="Load 1m" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="five" stroke="#82ca9d" name="Load 5m" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="fifteen" stroke="#ffc658" name="Load 15m" activeDot={{ r: 8 }} />
        </LineChart>
      </ResponsiveContainer>
    );
  };
  