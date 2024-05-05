import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Card, Spin, Typography } from 'antd';
const { Paragraph } = Typography;

export const LogViewer = () => {
    const [logs, setLogs] = useState('');
    const [loading, setLoading] = useState(true);
  
    useEffect(() => {
      const fetchLogs = async () => {
        const token = localStorage.getItem("token");
        if (!token) {
          console.log("Token not found, redirecting to login.");
          // 实际应用中应重定向到登录页面
          return;
        }
        try {
          const response = await axios.get("http://localhost:8080/logger/get", {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
          if (response.data.code === 200) {
            setLogs(response.data.result);
          } else {
            console.error("Failed to fetch logs:", response.data.message);
          }
        } catch (error) {
          console.error("Error fetching logs:", error);
        }
        setLoading(false);
      };
  
      fetchLogs();
    }, []);
  
    if (loading) {
      return <Spin size="large" />;
    }
  
    // 处理日志字符串，将其转换为可读格式
    const formattedLogs = logs.split('\n').map((log, index) => (
      <Paragraph key={index}>
        {log}
      </Paragraph>
    ));
  
    return (
      <Card title="日志查看" bordered={false} style={{ margin: '20px' }}>
        {formattedLogs.length > 0 ? formattedLogs : <Paragraph>没有日志数据</Paragraph>}
      </Card>
    );
  };
  