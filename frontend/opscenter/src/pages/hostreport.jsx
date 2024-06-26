import React, { useState, useEffect,useCallback } from 'react';
import axios from 'axios';
import { useSearchParams } from 'react-router-dom';
import { Spin } from 'antd';
import moment from 'moment';
import { CPUGraph } from '../components/cpugraph';
import { MemoryGraph } from '../components/memorygraph';
import { DiskGraph } from '../components/diskgraph';
import { LoadGraph } from '../components/loadgraph';

const processData = (rawData) => {
  const groupedData = rawData.reduce((acc, item) => {
    const time = moment(item._time).format('HH:mm:ss');
    if (!acc[time]) {
      acc[time] = { time };
    }
    acc[time][item._field] = item._value;
    return acc;
  }, {});

  return Object.values(groupedData);
};

export const HostReport = () => {
  const [searchParams] = useSearchParams();
  const hostId = searchParams.get('host_id');
  const [cpuData, setCpuData] = useState([]);
  const [memoryData, setMemoryData] = useState([]);
  const [diskData, setDiskData] = useState([]);
  const [loadData, setLoadData] = useState([]);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async (mt, setData) => {
    try {
      const response = await axios.get(`http://localhost:8080/reporter/report?hostId=${hostId}&mt=${mt}`);
      const formattedData = processData(response.data.result);
      setData(formattedData);
    } catch (error) {
      console.error('Error fetching report:', error);
    }
  }, [hostId]);

  useEffect(() => {
    const fetchAllData = async () => {
      await fetchData('cpu', setCpuData);
      await fetchData('memory', setMemoryData);
      await fetchData('disk', setDiskData);
      await fetchData('load', setLoadData);
    };

    fetchAllData(); // Initial fetch

    const intervalId = setInterval(() => {
      fetchAllData();
    }, 3000); // Fetch data every 3 seconds

    return () => clearInterval(intervalId); // Cleanup interval on component unmount
  }, [fetchData]);

  useEffect(() => {
    setLoading(false);
  }, [cpuData, memoryData, diskData, loadData]); // Update loading state when any data updates

  if (loading) return <Spin size="large" />;

  return (
    <div>
      <h3>CPU指标</h3>
      <CPUGraph data={cpuData} />
      <h3>内存指标</h3>
      <MemoryGraph data={memoryData} />
      <h3>磁盘指标</h3>
      <DiskGraph data={diskData} />
      <h3>负载指标</h3>
      <LoadGraph data={loadData} />
    </div>
  );
};
