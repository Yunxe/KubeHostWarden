import {
  EyeOutlined,
  FileOutlined,
  LaptopOutlined,
  NotificationOutlined,
  PlusOutlined,
} from "@ant-design/icons";
import { Layout, Menu } from "antd";
import Sider from "antd/es/layout/Sider";
import React from "react";
import { Link } from 'react-router-dom';

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
              <Link to="/hosts/view">查看主机</Link>
            </Menu.Item>
          </SubMenu>
          <Menu.Item key="3" icon={<NotificationOutlined />}>
            <Link to="/alerts">警报</Link>
          </Menu.Item>
          <Menu.Item key="4" icon={<FileOutlined />}>
            <Link to="/logs">日志</Link>
          </Menu.Item>
          {/* 更多菜单项... */}
        </Menu>
      </Sider>
    </Layout>
  );
};
