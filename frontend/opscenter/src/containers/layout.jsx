import { Layout } from "antd";
import React from "react";
import { HeaderLayout } from "./header";
import { SiderLayout } from "./sider";
import { ContentLayout } from "./content";

export function LayoutComponent({ children }) {
  return (
    <Layout>
      <HeaderLayout />
      <SiderLayout />
      <ContentLayout>{children}</ContentLayout>
    </Layout>
  );
}
