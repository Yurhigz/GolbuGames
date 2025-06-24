import React from "react";
import "./Input.css";

export const Input = ({ value, onChange, className = "", ...props }) => {
  return (
    <input
      className={`input-field ${className}`}
      value={value}
      onChange={onChange}
      {...props}
    />
  );
};
