import React, { useState } from "react";
import axios from "axios";
import { Form, Input, Button, Select ,Alert} from "antd";
import AlertModal from "../components/alertmodal";

export const HostAdd = () => {
  const [form] = Form.useForm();
  const [modalVisible, setModalVisible] = useState(false);
  const [modalMessage, setModalMessage] = useState("");

  const onFinish = async (values) => {
    const token = localStorage.getItem("token"); // 获取token
    if (!token) {
      window.location.href = "/login"; // 未登录或会话过期时跳转到登录页
      setModalMessage("失败");
      setModalVisible(true); // 显示模态框
      return; // 如果没有token，则终止请求
    }
    values.port = parseInt(values.port);
    try {
      const response = await axios.post(
        "http://localhost:8080/host/register",
        values,
        {
          headers: {
            Authorization: `Bearer ${token}`, // 将token添加到请求头中
          },
        }
      );
      if (response.data.code === 200) {
        setModalMessage("成功");
        setModalVisible(true);
      } else {
        setModalMessage(`失败: ${response.data.message}`);
        setModalVisible(true);
      }
    } catch (error) {
      setModalMessage("失败: ${error.response}");
      setModalVisible(true);
    }
  };

  const handleOk = () => {
    setModalVisible(false); // 关闭模态框
  };

  return (
    <div>
      <h2 style={{ textAlign: "center" }}>添加你的主机</h2>
      <Form
        form={form}
        name="hostForm"
        onFinish={onFinish}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 8 }}
      >
        <Form.Item
          label="ip地址"
          name="endpoint"
          rules={[{ required: true, message: "请输入主机地址！" }]}
        >
          <Input placeholder="e.g. 111.111.111.111" />
        </Form.Item>
        <Form.Item
          label="端口"
          name="port"
          rules={[{ required: true, message: "请输入端口号！" }]}
        >
          <Input type="text" placeholder="e.g. 22" />
        </Form.Item>
        <Form.Item
          label="用户名"
          name="user"
          rules={[{ required: true, message: "请输入用户名！" }]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="密码"
          name="password"
          rules={[{ required: true, message: "请输入密码！" }]}
        >
          <Input.Password />
        </Form.Item>
        <Form.Item
          label="操作系统类型"
          name="ostype"
          rules={[{ required: true, message: "请选择操作系统类型！" }]}
        >
          <Select>
            <Select.Option value="linux">Linux</Select.Option>
            <Select.Option value="darwin">MacOS</Select.Option>
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
