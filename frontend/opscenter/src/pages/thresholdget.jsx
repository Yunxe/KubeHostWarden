import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Table, Spin, Typography, Card, Button, Popconfirm, message } from 'antd';

const { Title } = Typography;

export const ThresholdGet = () => {
  const [thresholds, setThresholds] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchThresholds = async () => {
      const token = localStorage.getItem("token");
      if (!token) {
        console.log("Token not found, redirecting to login.");
        // 实际应用中应重定向到登录页面
        return;
      }
      try {
        const response = await axios.get("http://127.0.0.1:8080/alarm/getthreshold", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (response.data.code === 200) {
          setThresholds(response.data.result);
        } else {
          console.error("Failed to fetch thresholds:", response.data.message);
        }
      } catch (error) {
        console.error("Error fetching thresholds:", error);
      }
      setLoading(false);
    };

    fetchThresholds();
  }, []);

  const handleDelete = async (id) => {
    const token = localStorage.getItem("token");
    if (!token) {
      console.log("Token not found, redirecting to login.");
      // 实际应用中应重定向到登录页面
      return;
    }
    try {
      const response = await axios.post(
        "http://localhost:8080/alarm/deletethreshold",
        { id },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (response.data.code === 200) {
        message.success("阈值成功删除");
        // 更新阈值列表
        setThresholds(thresholds.filter((threshold) => threshold.id !== id));
      } else {
        message.error("Failed to delete threshold: " + response.data.message);
      }
    } catch (error) {
      console.error("Error deleting threshold:", error);
      message.error("Error deleting threshold");
    }
  };

  const columns = [
    {
      title: '指标',
      dataIndex: 'metric',
      key: 'metric',
    },
    {
      title: '子指标',
      dataIndex: 'sub_metric',
      key: 'sub_metric',
    },
    {
      title: '阈值',
      dataIndex: 'threshold',
      key: 'threshold',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Popconfirm
          title="确认删除？"
          onConfirm={() => handleDelete(record.id)}
          okText="是"
          cancelText="否"
        >
          <Button type="link" danger>
            删除
          </Button>
        </Popconfirm>
      ),
    },
  ];

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <Card title="查看阈值" bordered={false} style={{ margin: '20px' }}>
      <Title level={2}>阈值列表</Title>
      <Table columns={columns} dataSource={thresholds} rowKey="id" />
    </Card>
  );
};
