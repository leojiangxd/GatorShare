import { Eye, EyeOff, Lock, Mail, User } from "lucide-react";
import NavBar from "./components/NavBar";
import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";

const Register = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isFormValid, setIsFormValid] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleRegister = () => {
    alert(`Username: ${username}\nEmail: ${email}\nPassword: ${password}`);
  };

  useEffect(() => {
    const usernameValid = /^[A-Za-z0-9]{2,24}$/.test(username);
    const emailValid = /^[a-zA-Z0-9._]+@ufl\.edu$/.test(email);
    const passwordValid = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}$/.test(
      password
    );

    setIsFormValid(usernameValid && emailValid && passwordValid);
  }, [username, email, password]);

  return (
    <div className="min-h-screen flex flex-col bg-base-200">
      <NavBar />
      <div className="flex flex-1 items-center justify-center">
        <fieldset className="fieldset w-xs bg-base-100 border border-base-300 p-10 rounded-box">
          <legend className="fieldset-legend text-2xl">
            Create an Account
          </legend>

          <div>
            <label className="input w-full validator">
              <User className="opacity-50 h-[1em]" />
              <input
                type="text"
                required
                placeholder="Username"
                pattern="[A-Za-z0-9]*"
                minLength="3"
                maxLength="24"
                title="Must be alphanumeric"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </label>
            <p className="validator-hint hidden">
              Must be 3 to 24 characters containing only letters or numbers
            </p>
          </div>

          <div className="mt-2">
            <label className="input w-full validator">
              <Mail className="opacity-50 h-[1em]" />
              <input
                type="email"
                placeholder="example@ufl.edu"
                required
                pattern="^[a-zA-Z0-9._]+@ufl\.edu$"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </label>
            <div className="validator-hint hidden">
              Enter a valid UF email
            </div>
          </div>

          <div className="mt-2">
            <label className="input w-full validator">
              <Lock className="opacity-50 h-[1em]" />
              <input
                type={showPassword ? "text" : "password"}
                required
                placeholder="Password"
                minLength="8"
                pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}"
                title="Must be more than 8 characters, including number, lowercase letter, uppercase letter"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              {showPassword ? (
                <EyeOff
                  className="opacity-50 h-[1em] cursor-pointer"
                  onClick={() => setShowPassword(!showPassword)}
                />
              ) : (
                <Eye
                  className="opacity-50 h-[1em] cursor-pointer"
                  onClick={() => setShowPassword(!showPassword)}
                />
              )}
            </label>
            <p className="validator-hint hidden">
              Must be more than 8 characters, including
              <br />• At least one number
              <br />• At least one lowercase letter
              <br />• At least one uppercase letter
            </p>
          </div>

          <button
            className="btn btn-primary mt-4"
            onClick={handleRegister}
            disabled={!isFormValid}
          >
            Register
          </button>

          <Link to="/login" className="link link-primary">
            Have an account? Login here!
          </Link>
        </fieldset>
      </div>
    </div>
  );
};

export default Register;
