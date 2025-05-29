import React from 'react';
import { Link } from 'react-router-dom';
import { AiFillHome } from 'react-icons/ai';
import './SideNavUser.css';

const SideNavUser = () => {
    return (
        <nav className="sidenav-user">
            <Link to="/login" className="nav-item">
                <AiFillHome className="nav-icon" />
                <span>Login</span>
            </Link>
        </nav>
    );
};

export default SideNavUser;