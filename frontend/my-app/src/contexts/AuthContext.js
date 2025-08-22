// src/context/AuthContext.js
import { createContext, useState, useEffect } from "react";
import {redirect} from "react-router-dom";

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);

    // Charger le user depuis localStorage au démarrage
    useEffect(() => {
        const storedUser = localStorage.getItem("user");
        if (storedUser) {
            setUser(JSON.parse(storedUser));
        }
    }, []);

    const AuthLogin = (userData, token) => {
        setUser(userData);
        localStorage.setItem("user", JSON.stringify(userData));
        localStorage.setItem("token", token);
    };

    const AuthLogout = () => {
        setUser(null);
        localStorage.removeItem("user");
        localStorage.removeItem("token");
        redirect('/')
    };

    return (
        <AuthContext.Provider value={{ user, AuthLogin, AuthLogout }}>
            {children}
        </AuthContext.Provider>
    );
};
