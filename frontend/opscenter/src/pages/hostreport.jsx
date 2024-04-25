import React, { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import axios from "axios";
import { Spin } from "antd";
import { CPUGraph } from "../components/cpugraph";

export const HostReport = () => {
  const [searchParams] = useSearchParams();
  const hostId = searchParams.get("host_id");
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchHostReport = () => {
      const token = localStorage.getItem("token"); // 获取token
      if (!token) {
        window.location.href = "/login"; // 未登录或会话过期时跳转到登录页
        return; // 如果没有token，则终止请求
      }
      axios
        .get(`http://localhost:8080/reporter/report?hostId=${hostId}`, {
          headers: {
            Authorization: `Bearer ${token}`, // 将token添加到请求头中
          },
        })
        .then((response) => {
          console.log("Polling data:", response.data);
          setData(response.data);
          setLoading(false);
          setError(null);
        })
        .catch((err) => {
          console.error("Error fetching report:", err);
          setError("Failed to fetch data");

          setLoading(false);
        });
    };

    const intervalId = setInterval(fetchHostReport, 5000); // 设置轮询，每5秒请求一次

    // 清理函数，当组件卸载时清除定时器
    return () => clearInterval(intervalId);
  }, [hostId]);

  if (loading) {
    return <Spin size="large" />;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <h2>Host Metrics</h2>
      {loading ? <p>Loading...</p> : <CPUGraph data={data} />}
    </div>
  );
};
