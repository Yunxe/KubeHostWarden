import { Layout } from "antd";
import React from "react";
import { HeaderLayout } from "./header";
import { SiderLayout } from "./sider";
import { ContentLayout } from "./content";
import { Outlet } from "react-router-dom";

export function LayoutComponent() {
  return (
    <Layout>
      <HeaderLayout />
      <SiderLayout />
      <ContentLayout>
        <Outlet />  
      </ContentLayout>
    </Layout>
  );
}
