import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

const NavBar = () => {
  const isLoggedIn = true;
  const [searchTerm, setSearchTerm] = useState("");
  const navigate = useNavigate();

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleKeyDown = (e) => {
    if (e.key === "Enter" && searchTerm.trim() !== "") {
      navigate(`/search/${searchTerm}`);
    }
  };

  const handleLogout = () => {
    alert("Logged out!");
  };

  return (
    <div className="navbar bg-base-100 shadow-sm gap-2">
      <div className="flex flex-2">
        <Link to="/" className="btn btn-ghost text-xl">
          GatorShare
        </Link>
      </div>

      <div className="flex flex-3 justify-center">
        <input
          type="text"
          placeholder="Search"
          className="input input-bordered w-full flex-grow"
          value={searchTerm}
          onChange={handleSearchChange}
          onKeyDown={handleKeyDown}
        />
      </div>

      <div className="flex flex-2 justify-end">
        {isLoggedIn ? (
          <div className="flex gap-2">
          <Link to="/create" className="btn btn-ghost text-l">
            Create
          </Link>
            <div className="dropdown dropdown-end">
              <div
                tabIndex={0}
                role="button"
                className="btn btn-ghost btn-circle avatar"
              >
                <div className="w-10 rounded-full">
                  <img
                    alt="User Avatar"
                    src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"
                  />
                </div>
              </div>
              <ul
                tabIndex={0}
                className="menu menu-sm dropdown-content bg-base-200 rounded-box z-10 mt-3 w-52 p-2 shadow"
              >
                <li>
                  <Link to="/notifications" className="justify-between">
                    Notifications
                    <span className="badge bg-primary text-xs text-primary-content border-none">
                      New
                    </span>
                  </Link>
                </li>
                <li>
                  <Link to="/user/User1">Profile</Link>
                </li>
                <li>
                  <Link onClick={handleLogout}>Logout</Link>
                </li>
              </ul>
            </div>
          </div>
        ) : (
          <Link to="/login" className="btn btn-primary text-l">
            Login
          </Link>
        )}
      </div>
    </div>
  );
};

export default NavBar;
