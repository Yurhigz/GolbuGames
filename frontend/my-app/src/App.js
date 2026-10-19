import React, { useContext } from "react";
import { BrowserRouter as Router, Navigate, Route, Routes } from "react-router-dom";
import SideNav from "./components/SideNav";
import Footer from "./components/footer";
import Home from "./pages/homepage";
import Solo from "./pages/Solo";
import Multi from "./pages/Multi";
import Invite from "./pages/Invite";
import Tournament from "./pages/Tournament";
import Leaderboard from "./pages/Leaderboard";
import Login from "./pages/User/Login";
import Register from "./pages/User/Register";
import Profile from "./pages/Profile";
import About from "./pages/legal/About";
import Help from "./pages/legal/Help";
import FAQ from "./pages/legal/FAQ";
import Terms from "./pages/legal/Terms";
import Privacy from "./pages/legal/Privacy";
import Friends from "./pages/Friends";
import { AuthContext, AuthProvider } from "./contexts/AuthContext";
import { ThemeProvider } from "./contexts/ThemeContext";

const RequireAuth = ({ children }) => {
  const { user, isLoading } = useContext(AuthContext);

  if (isLoading) {
    return <p role="status">Chargement du compte…</p>;
  }

  return user ? children : <Navigate to="/login" replace />;
};

function App() {
  return (
    <ThemeProvider>
      <Router>
        <div className="app">
          <AuthProvider>
            <SideNav />
            <div className="main-content">
              <Routes>
                <Route path="/multiplayer/invite/:id" element={<Invite />} />
                <Route path="/" element={<Home />} />
                <Route path="/solo" element={<Solo />} />
                <Route path="/multi" element={<RequireAuth><Multi /></RequireAuth>} />
                <Route path="/tournament" element={<Tournament />} />
                <Route path="/leaderboard" element={<Leaderboard />} />
                <Route path="/friends" element={<RequireAuth><Friends /></RequireAuth>} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/profile" element={<RequireAuth><Profile /></RequireAuth>} />
                <Route path="/about" element={<About />} />
                <Route path="/help" element={<Help />} />
                <Route path="/faq" element={<FAQ />} />
                <Route path="/terms" element={<Terms />} />
                <Route path="/privacy" element={<Privacy />} />
              </Routes>
            </div>
          </AuthProvider>
          <Footer />
        </div>
      </Router>
    </ThemeProvider>
  );
}

export default App;