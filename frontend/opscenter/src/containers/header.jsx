import React from "react";
import { Layout, Menu, Button } from "antd";
import { UserOutlined } from "@ant-design/icons";

const { Header } = Layout;

export function HeaderLayout() {
  return (
    <Header
      style={{
        position: "fixed",
        zIndex: 1,
        width: "100%",
        padding: "0px 20px 20px 20px",
      }}
    >
      <div className="logo" />
      <Menu theme="dark" mode="horizontal" style={{ float: "left" }}>
        {/* 菜单项... */}
      </Menu>
      <div
        style={{
          float: "left",
          fontSize: "33px",
          color: "white",
          marginLeft: 0,
        }}
      >
        监控平台
      </div>
      <div style={{ float: "right" }}>
        <Button
          icon={<UserOutlined />}
          style={{ marginLeft: 10 }}
          onClick={() => {
            console.log("进入用户信息页面");
          }}
        ></Button>
      </div>
    </Header>
  );
}
