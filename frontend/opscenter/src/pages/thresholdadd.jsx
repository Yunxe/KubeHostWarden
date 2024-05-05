import React, { useState } from "react";
import axios from "axios";
import { Form, Input, Button, Select } from "antd";
import AlertModal from "../components/alertmodal";

export const ThresholdSetting = () => {
  const [form] = Form.useForm();
  const [modalVisible, setModalVisible] = useState(false);
  const [modalMessage, setModalMessage] = useState("");

  const onFinish = async (values) => {
    const token = localStorage.getItem("token"); // 获取token
    if (!token) {
      window.location.href = "/login"; // 未登录或会话过期时跳转到登录页
      setModalMessage("Authentication failed. Please login.");
      setModalVisible(true);
      return;
    }

    try {
      const response = await axios.post(
        "http://localhost:8080/thresholds/set",
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
  };

  const handleOk = () => {
    setModalVisible(false); // 关闭模态框
  };

  return (
    <div>
      <h2 style={{ textAlign: "center" }}>设置阈值</h2>
      <Form
        form={form}
        name="thresholdForm"
        onFinish={onFinish}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 8 }}
      >
        <Form.Item
          label="Host ID"
          name="host_id"
          rules={[{ required: true, message: "Please input the host ID!" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Metric"
          name="metric"
          rules={[{ required: true, message: "Please select a metric!" }]}
        >
          <Select>
            <Select.Option value="cpu">CPU</Select.Option>
            <Select.Option value="memory">Memory</Select.Option>
            <Select.Option value="disk">Disk</Select.Option>
            <Select.Option value="load">Load</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item
          label="Sub Metric"
          name="sub_metric"
          rules={[{ required: true, message: "Please select a sub metric!" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Threshold Value"
          name="threshold"
          rules={[{ required: true, message: "Please input the threshold value!" }]}
        >
          <Input type="number" />
        </Form.Item>
        <Form.Item
          label="Type"
          name="type"
          rules={[{ required: true, message: "Please select the type!" }]}
        >
          <Select>
            <Select.Option value="above">Above</Select.Option>
            <Select.Option value="below">Below</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item wrapperCol={{ offset: 11, span: 16 }}>
          <Button type="primary" htmlType="submit">
            提交
          </Button>
        </Form.Item>
      </Form>
      <AlertModal
        isVisible={modalVisible}
        handleOk={handleOk}
        message={modalMessage}
      />
    </div>
  );
};
