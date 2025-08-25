import { createContext, useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [token, setToken] = useState(null);
    const navigate = useNavigate();
    const logoutInProgress = useRef(false);

    useEffect(() => {
        const storedUser = localStorage.getItem("user");
        const storedToken = localStorage.getItem("token");

        if (storedUser) setUser(JSON.parse(storedUser));
        if (storedToken) setToken(storedToken);
    }, []);

    const AuthLogin = (userData, newToken) => {
        setUser(userData);
        setToken(newToken);

        localStorage.setItem("user", JSON.stringify(userData));
        localStorage.setItem("token", newToken);
    };


    const AuthLogout = async () => {
        if (logoutInProgress.current) return;
        logoutInProgress.current = true;

        setUser(null);
        setToken(null);
        localStorage.removeItem("user");
        localStorage.removeItem("token");

        try {
            await axios.post(`${process.env.REACT_APP_BACKEND_URL}/logout`, {}, { withCredentials: true });
        } catch (err) {
            console.error("Logout error:", err);
        }

        navigate("/login");
        logoutInProgress.current = false;
    };

    const getAuthHeaders = () => ({
        "Authorization": token ? `Bearer ${token}` : "",
        "Content-Type": "application/json"
    });

    const refreshJwtToken = async () => {
        try {
            const res = await axios.post(
                `${process.env.REACT_APP_BACKEND_URL}/refresh_token`,
                {},
                { withCredentials: true }
            );

            const data = res.data;

            if (data.access_token) {
                setToken(data.access_token);
                localStorage.setItem("token", data.access_token);
                return data.access_token;
            } else {
                AuthLogout();
                return null;
            }
        } catch (err) {
            console.error("Erreur refresh token :", err);
            AuthLogout();
            return null;
        }
    };

    return (
        <AuthContext.Provider value={{
            user,
            token,
            setToken,
            AuthLogin,
            AuthLogout,
            getAuthHeaders,
            refreshJwtToken
        }}>
            {children}
        </AuthContext.Provider>
    );
};
