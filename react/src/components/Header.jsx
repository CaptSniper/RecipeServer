import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import './Header.css';

export default function Header() {
  const [menuOpen, setMenuOpen] = useState(false);
  const location = useLocation();

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  const closeMenu = () => {
    setMenuOpen(false);
  };

  const isActive = (path) => {
    return location.pathname === path;
  };

  return (
    <header className="header">
      <div className="header-container">
        <Link to="/" className="header-logo" onClick={closeMenu}>
          Digital Cookbook
        </Link>

        {/* Desktop Navigation */}
        <nav className="header-nav desktop-nav">
          <Link 
            to="/" 
            className={`nav-link ${isActive('/') || isActive('/recipes') ? 'active' : ''}`}
          >
            Recipes
          </Link>
          <Link 
            to="/new" 
            className={`nav-link ${isActive('/new') ? 'active' : ''}`}
          >
            New Recipe
          </Link>
          <Link 
            to="/scrape" 
            className={`nav-link ${isActive('/scrape') ? 'active' : ''}`}
          >
            Scrape Recipe
          </Link>
        </nav>

        {/* Mobile Hamburger */}
        <button 
          className="hamburger" 
          onClick={toggleMenu}
          aria-label="Toggle menu"
        >
          <span className={`hamburger-line ${menuOpen ? 'open' : ''}`}></span>
          <span className={`hamburger-line ${menuOpen ? 'open' : ''}`}></span>
          <span className={`hamburger-line ${menuOpen ? 'open' : ''}`}></span>
        </button>
      </div>

      {/* Mobile Menu Dropdown */}
      <nav className={`mobile-nav ${menuOpen ? 'open' : ''}`}>
        <Link 
          to="/" 
          className={`nav-link ${isActive('/') || isActive('/recipes') ? 'active' : ''}`}
          onClick={closeMenu}
        >
          Recipes
        </Link>
        <Link 
          to="/new" 
          className={`nav-link ${isActive('/new') ? 'active' : ''}`}
          onClick={closeMenu}
        >
          New Recipe
        </Link>
        <Link 
          to="/scrape" 
          className={`nav-link ${isActive('/scrape') ? 'active' : ''}`}
          onClick={closeMenu}
        >
          Scrape Recipe
        </Link>
      </nav>
    </header>
  );
}

