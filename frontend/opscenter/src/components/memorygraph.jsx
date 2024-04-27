import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import moment from 'moment';

export const MemoryGraph = ({ data }) => {
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
          <Line type="monotone" dataKey="used" stroke="#8884d8" name="Used Memory" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="wired" stroke="#82ca9d" name="Wired Memory" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="unused" stroke="#ffc658" name="Unused Memory" activeDot={{ r: 8 }} />
          <Line type="monotone" dataKey="compressed" stroke="#ff7300" name="Compressed Memory" activeDot={{ r: 8 }} />
        </LineChart>
      </ResponsiveContainer>
    );
  };
  