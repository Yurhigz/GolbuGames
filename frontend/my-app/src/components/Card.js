import React from "react";
import "./Card.css";

export const Card = ({ children, className = "" }) => {
  return (
    <div className={`card-container rounded-2xl bg-white shadow-md ${className}`}>
      {children}
    </div>
  );
};

export const CardContent = ({ children, className = "" }) => {
  return (
    <div className={`card-content p-4 ${className}`}>
      {children}
    </div>
  );
};

