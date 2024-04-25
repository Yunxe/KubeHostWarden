import React, { useState, useEffect } from "react";
import { Card, Col, Row, Spin } from "antd";
import axios from "axios";
import { useNavigate } from 'react-router-dom'; 

export const HostRetrieve = () => {
  const [hosts, setHosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate(); 

  useEffect(() => {
    const fetchHosts = async () => {
      const token = localStorage.getItem("token"); // 获取token
      if (!token) {
        window.location.href = "/login"; // 未登录或会话过期时跳转到登录页
        return; // 如果没有token，则终止请求
      }
      try {
        const response = await axios.get(
          "http://localhost:8080/reporter/retrieve",
          {
            headers: {
              Authorization: `Bearer ${token}`, // 将token添加到请求头中
            },
          }
        );
        if (response.data.code === 200) {
          setHosts(response.data.result);
        } else {
          console.error("Failed to retrieve hosts:", response.data.message);
        }
      } catch (error) {
        console.error("Error fetching hosts:", error);
      }
      setLoading(false);
    };

    fetchHosts();
  }, []);

  const handleCardClick = (hostId) => {
    navigate(`/hosts/report?host_id=${hostId}`); // 使用 navigate 进行路由跳转
  };

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <div style={{ padding: "30px" }}>
      <Row gutter={16}>
        {hosts.map((host) => (
          <Col key={host.id} span={8}>
            <Card
              title={host.hostname}
              bordered={false}
              hoverable
              onClick={() => handleCardClick(host.id)} // 添加点击事件处理器
              style={{ cursor: "pointer" }}
            >
              <p>
                OS: {host.os} {host.os_version}
              </p>
              <p>
                Kernel: {host.kernel} {host.kernel_version}
              </p>
              <p>Architecture: {host.arch}</p>
              <p>IP Address: {host.ip_addr}</p>
              <p>Memory Total: {host.memory_total}</p>
              <p>Disk Total: {host.disk_total}</p>
              <p>Owner: {host.owner}</p>
              <p>Created At: {new Date(host.created_at).toLocaleString()}</p>
              <p>Updated At: {new Date(host.updated_at).toLocaleString()}</p>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  );
};
