import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { LayoutComponent } from "./containers/layout";
import { HostAdd } from "./pages/hostadd";
import { HostRetrieve } from "./pages/hostretrieve";
import { HostReport } from "./pages/hostreport";
import Login from "./components/login";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<LayoutComponent />}>
          {/* <Route index element={<Home />} /> */}
          <Route path="hosts/add" element={<HostAdd />} />
          <Route path="hosts/retrieve" element={<HostRetrieve />} />
          <Route path="hosts/report" element={<HostReport />} />
          {/* 在这里可以继续添加更多子路由 */}
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
