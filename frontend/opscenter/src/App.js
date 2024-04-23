import React from "react";
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { LayoutComponent } from "./containers/layout";
import { Hosts } from "./pages/Hosts";

function App() {
  return (
    <Router>
      <LayoutComponent>
        <Routes>
          <Route path="/hosts/add" element={<Hosts />} />
          {/* 添加更多路由规则 */}
        </Routes>
      </LayoutComponent>
    </Router>
  );
}

export default App;
