
import { useContext, useRef, useState } from "react";

import { Input } from "../../components/Input";
import { Button } from "../../components/Button";
import { Link } from "react-router-dom";
import "./Register.css";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../../contexts/AuthContext";
import { API_BASE_URL } from "../../config/api";

const validatePassword = (password) => {
    const hasMinLength = password.length >= 6;
    const hasUppercase = /[A-Z]/.test(password);
    const hasNumber = /\d/.test(password);
    return hasMinLength && hasUppercase && hasNumber;
};


const Register = () => {
    const navigate = useNavigate();
    const loginRef = useRef();
    const passwordRef = useRef();

    const confirmPasswordRef = useRef();
    const [error, setError] = useState("");
    const { AuthLogin } = useContext(AuthContext);

    const handleRegister = async (event) => {
        event.preventDefault();
        const login = loginRef.current.value.trim();
        const password = passwordRef.current.value;
        const confirmPassword = confirmPasswordRef.current.value;

        if (!login || !password || !confirmPassword) {
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

        if (password !== confirmPassword) {
            setError("Les mots de passe ne correspondent pas.");
            return;
        }

        setError("");
        try {
            const response = await axios.post(`${API_BASE_URL}/create_user`, {
                Username: login,
                Password: password,
                Accountname: login,
            });
            const { access_token } = response.data;
            AuthLogin({ id: response.data.user_id, username: login }, access_token);
            navigate("/");
        } catch {
            setError("Erreur lors de la requête");
        }

    };

    return (
        <div className="register-container">
            <h2 className="register-title">Créer un compte</h2>
            <form onSubmit={handleRegister}>
                <div className="input-group">
                    <Input
                        type="text"
                        placeholder="Nom d'utilisateur"
                        autoComplete="username"
                        maxLength={20}
                        ref={loginRef}
                    />
                </div>
                <div className="input-group">
                    <Input
                        type="password"
                        placeholder="Mot de passe"
                        autoComplete="new-password"
                        maxLength={20}
                        ref={passwordRef}
                    />
                </div>
                <div className="input-group">
                    <Input
                        type="password"
                        placeholder="Confirmer le mot de passe"
                        autoComplete="new-password"
                        maxLength={20}
                        ref={confirmPasswordRef}
                    />
                </div>
                {error && <p className="form-error" role="alert">{error}</p>}
                <Button type="submit">S'inscrire</Button>
            </form>
            <div className="login-link">
                <p>Déjà un compte ? <Link to="/login">Se connecter</Link></p>
            </div>
        </div>
    );
};

export default Register;
