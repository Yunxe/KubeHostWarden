import React, { useState } from "react";
import { Form, Input, Button, Checkbox, Card } from "antd";
import { UserOutlined, LockOutlined, MailOutlined } from "@ant-design/icons";
import axios from "axios";

const Login = () => {
  const [loading, setLoading] = useState(false);
  const [mode, setMode] = useState("login");

  const onFinish = async (values) => {
    setLoading(true);
    const endpoint = mode === "login" ? "http://localhost:8080/user/login" : "http://localhost:8080/user/register";
    const data = mode === "login" ? { email: values.email, password: values.password } :
      { username: values.username, password: values.password, email: values.email };

    try {
      const response = await axios.post(endpoint, data);
      const { code, result } = response.data;
      if (code === 200) {
        if (mode === "login") {
          localStorage.setItem("token", result.token);
          window.location.href = "/";  // 只有在登录成功时跳转到根路由
        } else {
          setMode("login");  // 注册成功，切换到登录模式
          console.log("注册成功，请登录");
        }
      } else {
        console.error(mode === "login" ? "登录失败" : "注册失败", result.message);
      }
    } catch (error) {
      console.error(mode === "login" ? "登录失败:" : "注册失败:", error);
    } finally {
      setLoading(false);
    }
  };

  const toggleMode = () => {
    setMode(prevMode => prevMode === "login" ? "register" : "login");
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <Card title={mode === "login" ? "登录" : "注册"} style={{ width: 300 }}>
        <Form
          name="login_form"
          initialValues={{ remember: true }}
          onFinish={onFinish}
        >
          {mode === "register" && (
            <Form.Item
              name="username"
              rules={[{ required: true, message: "请输入您的用户名!" }]}
            >
              <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="用户名" />
            </Form.Item>
          )}
          <Form.Item
            name="email"
            rules={[{ required: true, message: "请输入您的邮箱!" }]}
          >
            <Input prefix={<MailOutlined className="site-form-item-icon" />} placeholder="邮箱" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: "请输入您的密码!" }]}
          >
            <Input
              prefix={<LockOutlined className="site-form-item-icon" />}
              type="password"
              placeholder="密码"
            />
          </Form.Item>
          {mode === "login" && (
            <Form.Item>
              <Form.Item name="remember" valuePropName="checked" noStyle>
                <Checkbox>记住我</Checkbox>
              </Form.Item>
            </Form.Item>
          )}
          <Form.Item>
            <Button type="primary" htmlType="submit" className="login-form-button" loading={loading}>
              {mode === "login" ? "登录" : "注册"}
            </Button>
            <Button type="link" onClick={toggleMode} style={{ float: "right" }}>
              {mode === "login" ? "注册账号" : "返回登录"}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Login;
