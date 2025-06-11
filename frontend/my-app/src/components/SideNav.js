import React from 'react';
import { Link } from 'react-router-dom';
import { AiFillHome } from 'react-icons/ai';
import { GiPerspectiveDiceSixFacesRandom } from 'react-icons/gi';
import { IoGameController } from 'react-icons/io5';
import { FaTrophy, FaChartLine, FaUserPlus } from 'react-icons/fa';
import './SideNav.css';

const SideNav = () => {
    return (
        <nav className="sidenav">
            <Link to="/" className="nav-item">
                <AiFillHome className="nav-icon" />
                <span>Accueil</span>
            </Link>
            <Link to="/solo" className="nav-item">
                <GiPerspectiveDiceSixFacesRandom className="nav-icon" />
                <span>Jeu Solo</span>
            </Link>
            <Link to="/multi" className="nav-item">
                <IoGameController className="nav-icon" />
                <span>Jeu Multi</span>
            </Link>
            <Link to="/tournament" className="nav-item">
                <FaTrophy className="nav-icon" />
                <span>Tournoi</span>
            </Link>
            <Link to="/leaderboard" className="nav-item">
                <FaChartLine className="nav-icon" />
                <span>Classement</span>
            </Link>
            <Link to="/friends" className="nav-item">
                <FaUserPlus className="nav-icon" />
                <span>Friends</span>
            </Link>
        </nav>
    );
};

export default SideNav;