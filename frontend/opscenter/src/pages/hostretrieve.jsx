import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import {
  Card,
  Col,
  Row,
  Spin,
  Menu,
  Dropdown,
  Button,
  Modal,
  Form,
  Input,
  Select,
} from "antd";
import { MoreOutlined } from "@ant-design/icons"; // 引入更多操作图标

export const HostRetrieve = () => {
  const [hosts, setHosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [currentHostId, setCurrentHostId] = useState(null);
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [modalVisible, setModalVisible] = useState(false);
  const [modalMessage, setModalMessage] = useState("");

  const showModal = (hostId) => {
    setCurrentHostId(hostId);
    setIsModalVisible(true);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
  };

  const handleMenuClick = (hostId, key) => {
    if (key === "setThreshold") {
      showModal(hostId);
    } else if (key === "deleteHost") {
      if (window.confirm("您确定要删除此主机吗？")) {
        deleteHost(hostId);
      }
    }
  };

  const deleteHost = async (hostId) => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        window.location.href = "/login";
        return;
      }
      const response = await axios.post(
        `http://localhost:8080/host/delete`,
        { hostId: hostId },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (response.data.code === 200) {
        alert("删除主机成功！");
        window.location.reload();
      } else {
        alert("删除主机失败！");
      }
    } catch (e) {
      console.error("删除主机失败：", e);
    }
  };

  const onFinish = async (values) => {
    const token = localStorage.getItem("token"); // 获取token
    if (!token) {
      window.location.href = "/login"; // 未登录或会话过期时跳转到登录页
      setModalMessage("Authentication failed. Please login.");
      setModalVisible(true);
      return;
    }
    values.threshold = parseFloat(values.threshold);
    values.host_id = currentHostId;

    try {
      const response = await axios.post(
        "http://localhost:8080/alarm/setthreshold",
        values,
        {
          headers: {
            Authorization: `Bearer ${token}`, // 将token添加到请求头中
          },
        }
      );
      if (response.data.code === 200) {
        setModalMessage("Threshold set successfully");
        setModalVisible(true);
      } else {
        setModalMessage(`Failed to set threshold: ${response.data.message}`);
        setModalVisible(true);
      }
    } catch (error) {
      setModalMessage(`Failed to set threshold: ${error.toString()}`);
      setModalVisible(true);
    }
    // 在这里添加调用后端API的代码
  };

  const menu = (hostId) => (
    <Menu
      onClick={(e) => {
        e.domEvent.stopPropagation();
        handleMenuClick(hostId, e.key);
      }}
    >
      <Menu.Item key="setThreshold">设置阈值</Menu.Item>
      <Menu.Item key="deleteHost">删除主机</Menu.Item>
    </Menu>
  );

  useEffect(() => {
    const fetchHosts = async () => {
      const token = localStorage.getItem("token");
      if (!token) {
        window.location.href = "/login";
        return;
      }
      try {
        const response = await axios.get(
          "http://localhost:8080/host/retrieve",
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        setHosts(response.data.result);
      } catch (error) {
        console.error("获取主机数据失败:", error);
      }
      setLoading(false);
    };

    fetchHosts();
  }, []);

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <div style={{ padding: "30px" }}>
      <Row gutter={16}>
        {hosts.map((host) => (
          <Col key={host.id} span={8}>
            <Card
              title={
                <div
                  style={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                  }}
                >
                  {host.hostname}
                  <Dropdown overlay={menu(host.id)} trigger={["click"]}>
                    <Button
                      type="text"
                      icon={<MoreOutlined />}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </Dropdown>
                </div>
              }
              bordered={false}
              hoverable
              onClick={() => navigate(`/hosts/report?host_id=${host.id}`)}
              style={{ cursor: "pointer" }}
            >
              <p>
                操作系统: {host.os} {host.os_version}
              </p>
              <p>
                内核版本: {host.kernel} {host.kernel_version}
              </p>
              <p>架构: {host.arch}</p>
              <p>IP地址: {host.ip_addr}</p>
              <p>总内存: {host.memory_total}</p>
              <p>总硬盘空间: {host.disk_total}</p>
              <p>所有者: {host.owner}</p>
              <p>创建时间: {new Date(host.created_at).toLocaleString()}</p>
              <p>更新时间: {new Date(host.updated_at).toLocaleString()}</p>
            </Card>
          </Col>
        ))}
      </Row>
      <Modal
        title="设置阈值"
        visible={isModalVisible}
        onCancel={handleCancel}
        footer={null}
      >
        <Form form={form} onFinish={onFinish}>
          <Form.Item
            name="metric"
            rules={[{ required: true, message: "请选择指标类型！" }]}
          >
            <Select placeholder="选择指标类型">
              <Select.Option value="cpu">CPU</Select.Option>
              <Select.Option value="memory">内存</Select.Option>
              <Select.Option value="disk">磁盘</Select.Option>
              <Select.Option value="load">负载</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="sub_metric"
            rules={[{ required: true, message: "请输入子指标！" }]}
          >
            <Input type="text" placeholder="输入子指标" />
          </Form.Item>
          <Form.Item
            name="threshold"
            rules={[{ required: true, message: "请输入阈值！" }]}
          >
            <Input type="number" placeholder="输入阈值" />
          </Form.Item>
          <Form.Item
            name="type"
            rules={[{ required: true, message: "请选择类型！" }]}
          >
            <Select placeholder="选择类型">
              <Select.Option value="above">高于</Select.Option>
              <Select.Option value="below">低于</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              设置阈值
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
