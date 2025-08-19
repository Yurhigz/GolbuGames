import { useRef } from "react";
import { Input } from "../../components/Input";
import { Button } from "../../components/Button";
import { Link } from "react-router-dom";
import "./Register.css";

const Register = () => {
    const loginRef = useRef();
    const passwordRef = useRef();

    const handleRegister = () => {
        const newUser = {
            login: loginRef.current.value,
            password: passwordRef.current.value,
        };
        console.log("Création de compte :", newUser);
    };

    return (
        <div className="register-container">
            <h2 className="register-title">Créer un compte</h2>
            <div className="input-group">
                <Input type="text" placeholder="Nom d'utilisateur" ref={loginRef} />
            </div>
            <div className="input-group">
                <Input type="password" placeholder="Mot de passe" ref={passwordRef} />
            </div>
            <Button onClick={handleRegister}>S'inscrire</Button>
            <div className="login-link">
                <p>Déjà un compte ? <Link to="/login">Se connecter</Link></p>
            </div>
        </div>
    );
};

export default Register;
