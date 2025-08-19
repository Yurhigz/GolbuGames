import { useRef, useState } from "react";
import { Input } from "../../components/Input";
import { Button } from "../../components/Button";
import { Link } from "react-router-dom";
import "./Register.css";
import axios from "axios";
import { useNavigate } from "react-router-dom";

// Fonction de nettoyage contre les injections (XSS)
const sanitizeInput = (str) => {
    const temp = document.createElement("div");
    temp.textContent = str;
    return temp.innerHTML;
};

// Vérification du mot de passe : 6+ caractères, 1 majuscule, 1 chiffre
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

    const handleRegister = async () => {
        const rawLogin = loginRef.current.value;
        const rawPassword = passwordRef.current.value;
        const rawConfirmPassword = confirmPasswordRef.current.value;

        const login = sanitizeInput(rawLogin.trim());
        const password = sanitizeInput(rawPassword.trim());
        const confirmPassword = sanitizeInput(rawConfirmPassword.trim());

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

        setError(""); // Réinitialiser les erreurs

        const newUser = {
            login,
            password,
        };

        try {
            const response = await axios.post("http://localhost:3001/create_user", {
                Username: login,
                Password: password,
                Accountname: login,
            });
            navigate("/");
            console.log("Réponse serveur :", response.data);
        } catch (error) {
            setError("Erreur lors de la requête");
            console.error("Erreur lors de la requête :", error);
        }
    };

    return (
        <div className="register-container">
            <h2 className="register-title">Créer un compte</h2>

            {error && <div className="register-error">{error}</div>}

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
            <div className="input-group">
                <Input
                    type="password"
                    placeholder="Confirmer le mot de passe"
                    ref={confirmPasswordRef}
                    autoComplete="off"
                />
            </div>

            <Button onClick={handleRegister}>S'inscrire</Button>

            <div className="login-link">
                <p>Déjà un compte ? <Link to="/login">Se connecter</Link></p>
            </div>
        </div>
    );
};

export default Register;
