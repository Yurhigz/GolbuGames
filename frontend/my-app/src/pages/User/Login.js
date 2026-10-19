
import { useRef, useState, useContext } from "react";
import { Input } from "../../components/Input";
import { Button } from "../../components/Button";
import { Link, useNavigate } from "react-router-dom";

import "./Login.css";
import axios from "axios";
import { AuthContext } from "../../contexts/AuthContext";
import { API_BASE_URL } from "../../config/api";

const Login = () => {
    const navigate = useNavigate();
    const loginRef = useRef();
    const passwordRef = useRef();

    const [error, setError] = useState("");
    const { AuthLogin } = useContext(AuthContext);

    const handleLogin = async (event) => {
        event.preventDefault();
        const login = loginRef.current.value.trim();
        const password = passwordRef.current.value;

        if (!login || !password) {
            setError("Veuillez remplir tous les champs.");
            return;
        }

        setError("");
        try {
            const response = await axios.post(`${API_BASE_URL}/login`, {
                Username: login,
                Password: password,
                Accountname: login,
            });

            const { access_token, user_id, username } = response.data;

            AuthLogin({ id: user_id, username }, access_token);
            navigate("/");
        } catch {
            setError("Login échoué !");
        }

    };

    return (
        <div className="login-container">
            <h2 className="login-title">Connexion</h2>
            <form onSubmit={handleLogin}>
                <div className="input-group">
                    <Input
                        type="text"
                        placeholder="Nom d'utilisateur"
                        autoComplete="username"
                        ref={loginRef}
                    />
                </div>
                <div className="input-group">
                    <Input
                        type="password"
                        placeholder="Mot de passe"
                        autoComplete="current-password"
                        ref={passwordRef}
                    />
                </div>
                {error && <p className="form-error" role="alert">{error}</p>}
                <Button type="submit">Connexion</Button>
            </form>
            <div className="register-link">
                <p>Pas encore de compte ? <Link to="/register">Créer un compte</Link></p>
            </div>
        </div>
    );
};

export default Login;
