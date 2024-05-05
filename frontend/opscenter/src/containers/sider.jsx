import {
  EyeOutlined,
  FileOutlined,
  LaptopOutlined,
  NotificationOutlined,
  AlertOutlined,
  PlusOutlined,
} from "@ant-design/icons";
import { Layout, Menu } from "antd";
import Sider from "antd/es/layout/Sider";
import React from "react";
import { Link } from "react-router-dom";

const { Header, Content, Footer } = Layout;
const { SubMenu } = Menu; // 引入 SubMenu 组件

export const SiderLayout = () => {
  return (
    <Layout style={{ marginTop: 64 }}>
      <Sider style={{ height: "100vh", position: "fixed", left: 0 }}>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={["1"]}
          style={{ height: "100%", borderRight: 0 }}
        >
          <SubMenu key="sub1" icon={<LaptopOutlined />} title="主机">
            <Menu.Item key="1" icon={<PlusOutlined />}>
              <Link to="/hosts/add">添加主机</Link>
            </Menu.Item>
            <Menu.Item key="2" icon={<EyeOutlined />}>
              <Link to="/hosts/retrieve">查看主机</Link>
            </Menu.Item>
          </SubMenu>
          <SubMenu key="sub2" icon={<AlertOutlined />} title="警报">
            <Menu.Item key="3" icon={<PlusOutlined />}>
              <Link to="/alarm/setthreshold">添加阈值</Link>
            </Menu.Item>
            <Menu.Item key="4" icon={<EyeOutlined />}>
              <Link to="/alarm/threshold">查看阈值</Link>
            </Menu.Item>
          </SubMenu>
          <SubMenu key="sub3" icon={<FileOutlined />} title="日志">
            <Menu.Item key="5" icon={<EyeOutlined />}>
              <Link to="/logger/view">查看日志</Link>
            </Menu.Item>
          </SubMenu>
          {/* 更多菜单项... */}
        </Menu>
      </Sider>
    </Layout>
  );
};
