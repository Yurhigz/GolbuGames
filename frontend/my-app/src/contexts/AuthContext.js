// src/context/AuthContext.js
import { createContext, useState, useEffect } from "react";

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    // Charger le user depuis localStorage au démarrage
    useEffect(() => {
        try {
            const storedUser = localStorage.getItem("user");
            if (storedUser) {
                const parsedUser = JSON.parse(storedUser);
                if (!parsedUser || typeof parsedUser !== "object" || !parsedUser.id) {
                    throw new Error("Stored user data is invalid");
                }
                setUser(parsedUser);
            }
        } catch (error) {
            console.error("Impossible de restaurer la session utilisateur :", error);
            localStorage.removeItem("user");
            localStorage.removeItem("token");
        } finally {
            setIsLoading(false);
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
    };

    return (
        <AuthContext.Provider value={{ user, isLoading, AuthLogin, AuthLogout }}>
            {children}
        </AuthContext.Provider>
    );
};
