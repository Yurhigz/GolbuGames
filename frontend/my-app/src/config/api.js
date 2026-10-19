const configuredApiURL = process.env.REACT_APP_API_URL || "http://localhost:3002";

export const API_BASE_URL = configuredApiURL.replace(/\/+$/, "");
export const WEBSOCKET_BASE_URL = (
    process.env.REACT_APP_WEBSOCKET_URL || API_BASE_URL.replace(/^http/, "ws")
).replace(/\/+$/, "");

export const authHeaders = () => ({
    Authorization: `Bearer ${localStorage.getItem("token") || ""}`,
});