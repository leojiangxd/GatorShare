import { Eye, EyeOff, Lock, User } from "lucide-react";
import NavBar from "./components/NavBar";
import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isFormValid, setIsFormValid] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const handleLogin = () => {
    alert(`Username: ${username}\nPassword: ${password}`);
  };

  useEffect(() => {
    setIsFormValid(username !== "" && password !== "");
  }, [username, password]);

  return (
    <div className="min-h-screen flex flex-col">
      <NavBar />
      <div className="flex flex-1 items-center justify-center">
        <fieldset className="fieldset w-xs bg-base-200 border border-base-300 p-10 rounded-box">
          <legend className="fieldset-legend text-2xl">Welcome Back</legend>

          <div>
            <label className="input w-full">
              <User className="opacity-50 h-[1em]" />
              <input
                type="text"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </label>
          </div>

          <div className="mt-2">
            <label className="input w-full">
              <Lock className="opacity-50 h-[1em]" />
              <input
                type={showPassword ? "text" : "password"}
                required
                placeholder="Password"
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
          </div>

          <button
            className="btn btn-primary mt-4"
            onClick={handleLogin}
            disabled={!isFormValid}
          >
            Login
          </button>

          <Link to="/register" className="link-info underline">
            Need an account? Register here!
          </Link>
        </fieldset>
      </div>
    </div>
  );
};

export default Login;
