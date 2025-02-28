import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getCsrfToken, getUsername } from "../../utils/functions";

const NavBar = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [loggedInUsername, setLoggedInUsername] = useState("");
  useEffect(() => {
    const fetchUsername = async () => {
      const username = await getUsername();
      setLoggedInUsername(username || "");
    };

    fetchUsername();
  }, []);
  const navigate = useNavigate();

  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  useEffect(() => {
    const checkUserLogin = async () => {
      try {
        const csrfToken = getCsrfToken();
        const response = await axios.get(`${apiBaseUrl}/api/v1/current-user`, {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        });

        if (response.data.username) {
          setIsLoggedIn(true);
        }
      } catch (error) {
        setIsLoggedIn(false);
      }
    };

    checkUserLogin();
  }, []);

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      if (searchTerm.trim() !== "") {
        navigate(`/search/${searchTerm}`);
      } else {
        navigate(`/`);
      }
    }
  };

  const handleLogout = async () => {
    const csrfToken = getCsrfToken();
    axios
      .post(
        `${apiBaseUrl}/api/v1/logout`,
        {},
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      )
      .then(() => {
        setIsLoggedIn(false);
      })
      .catch((error) => {
        console.alert("Logout failed:", error);
      });
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

      <div className="flex flex-2 justify-end items-center">
        {isLoggedIn ? (
          <div className="flex gap-2">
            <Link to="/create" className="btn btn-ghost text-l">
              Create
            </Link>
            <div className="dropdown dropdown-end">
              <div
                tabIndex={0}
                role="button"
                className="btn btn-primary btn-circle avatar bg-primary text-xl"
              >
                {loggedInUsername.trim().charAt(0).toUpperCase()}
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
                  <Link to={`/user/${loggedInUsername}`}>Profile</Link>
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
