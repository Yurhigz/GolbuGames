import { useRef } from "react";
import { Input } from "../../components/Input";
import { Button } from "../../components/Button";
import { Link } from "react-router-dom"; // Si tu utilises React Router
import "./Login.css";

const Login = () => {
    const loginRef = useRef();
    const passwordRef = useRef();

    const handleLogin = () => {
        const identifiants = {
            login: loginRef.current.value,
            password: passwordRef.current.value,
        };
        console.log("Connexion :", identifiants);
    };

    return (
        <div className="login-container">
            <h2 className="login-title">Connexion</h2>
            <div className="input-group">
                <Input type="text" placeholder="Nom d'utilisateur" ref={loginRef} />
            </div>
            <div className="input-group">
                <Input type="password" placeholder="Mot de passe" ref={passwordRef} />
            </div>
            <Button onClick={handleLogin}>Connexion</Button>
            <div className="register-link">
                <p>Pas encore de compte ? <Link to="/register">Créer un compte</Link></p>
            </div>
        </div>
    );
};

export default Login;
