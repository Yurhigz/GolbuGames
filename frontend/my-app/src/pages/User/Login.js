import { useRef, useState, useContext } from "react";
import { Input } from "../../components/Input";
import { Button } from "../../components/Button";
import {Link, useNavigate} from "react-router-dom";
import "./Login.css";
import axios from "axios";
import { AuthContext } from "../../contexts/AuthContext";

const sanitizeInput = (str) => {
    const temp = document.createElement("div");
    temp.textContent = str;
    return temp.innerHTML;
};

const validatePassword = (password) => {
    const hasMinLength = password.length >= 6;
    const hasUppercase = /[A-Z]/.test(password);
    const hasNumber = /\d/.test(password);

    return hasMinLength && hasUppercase && hasNumber;
};

const Login = () => {
    const navigate = useNavigate();
    const loginRef = useRef();
    const passwordRef = useRef();
    const [error, setError] = useState("");
    const { AuthLogin } = useContext(AuthContext);

    const handleLogin = async () => {
        const rawLogin = loginRef.current.value;
        const rawPassword = passwordRef.current.value;

        const login = sanitizeInput(rawLogin.trim());
        const password = sanitizeInput(rawPassword.trim());

        if (!login || !password) {
            setError("Veuillez remplir tous les champs.");
            return;
        }

        if (login.length > 20 || password.length > 20) {
            setError("Les champs ne doivent pas dépasser 20 caractères.");
            return;
        }

        if (!validatePassword(password)) {
            setError("Le mot de passe doit contenir au moins 6 caractères, une lettre majuscule et un chiffre.");
            return;
        }

        setError(""); // Réinitialise les erreurs

        const identifiants = {
            login,
            password,
        };

        console.log("Connexion sécurisée :", identifiants);
        try {
            const response = await axios.post("http://localhost:3001/login", {
                Username: login,
                Password: password,
                Accountname: login,
            });
            const { access_token } = response.data;

            AuthLogin({ login }, access_token);
            navigate("/");
        } catch (err) {
            console.error("Erreur login:", err);
            setError("Login échoué !");
        }
    };

    return (
        <div className="login-container">
            <h2 className="login-title">Connexion</h2>

            {error && <div className="login-error">{error}</div>}

            <div className="input-group">
                <Input
                    type="text"
                    placeholder="Nom d'utilisateur"
                    ref={loginRef}
                    autoComplete="off"
                />
            </div>
            <div className="input-group">
                <Input
                    type="password"
                    placeholder="Mot de passe"
                    ref={passwordRef}
                    autoComplete="off"
                />
            </div>

            <Button onClick={handleLogin}>Connexion</Button>

            <div className="register-link">
                <p>Pas encore de compte ? <Link to="/register">Créer un compte</Link></p>
            </div>
        </div>
    );
};

export default Login;
