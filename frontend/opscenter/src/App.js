import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { LayoutComponent } from "./containers/layout";
import { Hosts } from "./pages/Hosts";
import Login from "./components/login";

function App() {
 return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<LayoutComponent />}>
          {/* <Route index element={<Home />} /> */}
          <Route path="hosts/add" element={<Hosts />} />
          {/* 在这里可以继续添加更多子路由 */}
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
