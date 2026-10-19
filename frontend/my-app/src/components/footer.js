import React from 'react';
import { Link } from 'react-router-dom';
import './Footer.css';

const Footer = () => {
    return (
        <footer className="footer">
            <div className="footer-content">
                <div className="footer-section">
                    <h3>À propos</h3>
                    <Link to="/about">Qui sommes-nous</Link>
                    <a href="mailto:support@golbugames.com">Contact</a>
                </div>
                
                <div className="footer-section">
                    <h3>Aide & Support</h3>
                    <Link to="/help">Centre d'aide</Link>
                    <Link to="/faq">FAQ</Link>
                </div>
                
                <div className="footer-section">
                    <h3>Mentions légales</h3>
                    <Link to="/terms">Conditions d'utilisation</Link>
                    <Link to="/privacy">Politique de confidentialité</Link>
                </div>
            </div>
            <div className="footer-bottom">
                <p>&copy; {new Date().getFullYear()} GolbuGames. Tous droits réservés.</p>
            </div>
        </footer>
    );
};

export default Footer;