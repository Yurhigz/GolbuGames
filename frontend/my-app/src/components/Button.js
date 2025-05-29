import React from "react";
import "./Button.css";

export const Button = ({ children, onClick, type = "button", className = "", ...props }) => {
  return (
    <button
      type={type}
      className={`btn-green ${className}`}
      onClick={onClick}
      {...props}
    >
      {children}
    </button>
  );
};
