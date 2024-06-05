import { Layout, Menu } from "antd";
import React from "react";

const { Header, Content, Footer } = Layout;
const { SubMenu } = Menu; // 引入 SubMenu 组件

export const ContentLayout = ({ children }) => {
  return (
    <Layout style={{ marginLeft: 200 }}>
      <Content style={{ margin: "24px 16px 0", overflow: "initial" }}>
        <div className="site-layout-content">{children}</div> 
      </Content>
      <Footer style={{ textAlign: "center" }}>
      </Footer>
    </Layout>
  );
};
