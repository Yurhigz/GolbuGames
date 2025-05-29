import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import SideNav from "./components/SideNav";
import SideNavUser from "./components/SideNavUser";
import Footer from "./components/footer";
import Home from "./pages/homepage";
import Solo from "./pages/Solo";
import Multi from "./pages/Multi";
import Invite from "./pages/Invite";
import Tournament from "./pages/Tournament";
import Leaderboard from "./pages/Leaderboard";
import Login from "./pages/User/Login";
import Register from "./pages/User/Register";
import About from "./pages/legal/About";
import Help from "./pages/legal/Help";
import FAQ from "./pages/legal/FAQ";
import Terms from "./pages/legal/Terms";
import Privacy from "./pages/legal/Privacy";

function App() {
  return (
    <Router>
      <div className="app">
        <SideNav />
        <SideNavUser />
        <div className="main-content">
          <Routes>
            <Route path="/multiplayer/invite/:id" element={<Invite />} />
            <Route path="/" element={<Home />} />
            <Route path="/solo" element={<Solo />} />
            <Route path="/multi" element={<Multi />} />
            <Route path="/tournament" element={<Tournament />} />
            <Route path="/leaderboard" element={<Leaderboard />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/about" element={<About />} />
            <Route path="/help" element={<Help />} />
            <Route path="/faq" element={<FAQ />} />
            <Route path="/terms" element={<Terms />} />
            <Route path="/privacy" element={<Privacy />} />
          </Routes>
        </div>
        <Footer />
      </div>
    </Router>
  );
}

export default App;