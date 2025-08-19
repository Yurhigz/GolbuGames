import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { AiFillHome } from 'react-icons/ai';
import { GiPerspectiveDiceSixFacesRandom } from 'react-icons/gi';
import { IoGameController } from 'react-icons/io5';
import { FaTrophy, FaChartLine, FaUserPlus, FaBars } from 'react-icons/fa';
import './SideNav.css';

const SideNav = () => {
    const [isOpen, setIsOpen] = useState(false);

    const toggleSideNav = () => {
        setIsOpen(!isOpen);
    };

    const closeOnOutsideClick = (e) => {
        if (!e.target.closest('.sidenav') && !e.target.closest('.hamburger')) {
            setIsOpen(false);
        }
    };

    const handleLinkClick = () => {
        if (window.innerWidth < 1500) {
            setIsOpen(false);
        }
    };

    useEffect(() => {
        document.addEventListener('click', closeOnOutsideClick);
        return () => document.removeEventListener('click', closeOnOutsideClick);
    }, []);

    return (
        <>
            <button className="hamburger" onClick={toggleSideNav}>
                <FaBars />
            </button>
            <nav className={`sidenav ${isOpen ? 'open' : ''}`}>
                <Link to="/" className="nav-item" onClick={handleLinkClick}>
                    <AiFillHome className="nav-icon" />
                    <span>Accueil</span>
                </Link>
                <Link to="/solo" className="nav-item" onClick={handleLinkClick}>
                    <GiPerspectiveDiceSixFacesRandom className="nav-icon" />
                    <span>Jeu Solo</span>
                </Link>
                <Link to="/multi" className="nav-item" onClick={handleLinkClick}>
                    <IoGameController className="nav-icon" />
                    <span>Jeu Multi</span>
                </Link>
                <Link to="/tournament" className="nav-item" onClick={handleLinkClick}>
                    <FaTrophy className="nav-icon" />
                    <span>Tournoi</span>
                </Link>
                <Link to="/leaderboard" className="nav-item" onClick={handleLinkClick}>
                    <FaChartLine className="nav-icon" />
                    <span>Classement</span>
                </Link>
                <Link to="/friends" className="nav-item" onClick={handleLinkClick}>
                    <FaUserPlus className="nav-icon" />
                    <span>Amis</span>
                </Link>
                <Link to="/login" className="nav-item" onClick={handleLinkClick}>
                    <FaUserPlus className="nav-icon" />
                    <span>Connexion</span>
                </Link>
                <Link to="/profile" className="nav-item" onClick={handleLinkClick}>
                    <FaUserPlus className="nav-icon" />
                    <span>Profile</span>
                </Link>
            </nav>
        </>
    );
};

export default SideNav;
