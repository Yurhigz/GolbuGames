import { useContext, useCallback } from "react";
import { AuthContext } from "../contexts/AuthContext";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const REACT_BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

export const useRequest = () => {
    const { token, setToken, refreshJwtToken, AuthLogout } = useContext(AuthContext);
    const navigate = useNavigate();

    const sendRequest = useCallback(async (method, url, data = {}, jwt = false, cookie = false) => {
        const makeRequest = async (currentToken) => {
            const headers = jwt
                ? { "Authorization": `Bearer ${currentToken}`, "Content-Type": "application/json" }
                : { "Content-Type": "application/json" };

            const config = {
                method: method.toLowerCase(),
                url: `${REACT_BACKEND_URL}${url}`,
                headers,
                data,
                ...(cookie && { withCredentials: true })
            };

            if (method.toUpperCase() === "GET") {
                config.params = data;
                delete config.data;
            }

            return axios(config);
        };

        try {
            const response = await makeRequest(token);
            if (jwt && response.headers["x-new-token"]) setToken(response.headers["x-new-token"]);
            return response.data;
        } catch (err) {
            if (jwt && err.response && err.response.status === 401) {
                const newToken = await refreshJwtToken();
                if (newToken) {
                    try {
                        const retryResponse = await makeRequest(newToken);
                        if (retryResponse.headers["x-new-token"]) setToken(retryResponse.headers["x-new-token"]);
                        return retryResponse.data;
                    } catch {}
                }
            }
            AuthLogout();
            navigate("/login");
            throw err;
        }
    }, [token, setToken, refreshJwtToken, AuthLogout, navigate]);

    return { sendRequest };
};
