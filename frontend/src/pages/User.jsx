import React, { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";
import { ThumbsDown, ThumbsUp } from "lucide-react";
import { getCsrfToken, getUsername } from "../utils/functions";
import axios from "axios";

const Home = () => {
  const navigate = useNavigate();
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  
  const { id } = useParams();
  const [loggedInUsername, setLoggedInUsername] = useState("");
  const [ownProfile, setOwnProfile] = useState(false);

  // Toggle between view and edit modes.
  const [isEditing, setIsEditing] = useState(false);

  // Profile state.
  const [username, setUsername] = useState(id);
  const [email, setEmail] = useState("example@example.com");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [bio, setBio] = useState("");

  // New state for messaging
  const [isMessageBoxOpen, setIsMessageBoxOpen] = useState(false);
  const [messageContent, setMessageContent] = useState("");

  // Fetch Profile Details
  useEffect(() => {
    const fetchUserDetails = async () => {
      try {
        const loggedInUser = await getUsername();
        setLoggedInUsername(loggedInUser || "");

        const isOwnProfile = id.toLowerCase().trim() === loggedInUser.toLowerCase().trim();
        setOwnProfile(isOwnProfile);

        if (ownProfile) {
          // Fetch member details
          const memberResponse = await axios.get(`${apiBaseUrl}/api/v1/member/${id}`);
          const memberData = memberResponse.data.data;
          // Populate profile fields
          setEmail(memberData.email);
          setBio(memberData.bio);
        }
      } catch (error) {
        console.error("Error fetching user details:", error);
      }
    };

    fetchUserDetails();
  }, [apiBaseUrl, id, loggedInUsername]);

  // Loop through each post and calculate the total likes and dislikes
  const calcLikesAndDislikes = (posts) => {
    let totalLikes = 0;
    let totalDislikes = 0;

    posts.forEach((post) => {
      // Add the post's likes and dislikes
      totalLikes += post.likes;
      totalDislikes += post.dislikes;
    });
    return { totalLikes, totalDislikes };
  };

  // Handlers to simulate actions.
  const handleSave = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }

      const updatedInfo = {
        currentPassword,
        username,
        email,
        bio,
      };

      if (newPassword) {
        updatedInfo.newPassword = newPassword;
      }

      // Send update request
      await axios.put(`${apiBaseUrl}/api/v1/member`, updatedInfo, {
        headers: {
          "X-CSRF-Token": csrfToken,
        },
        withCredentials: true,
      });
      navigate(`/user/${username}`);
      setIsEditing(false);
    } catch (error) {
      console.error("Error updating profile:", error.response ? error.response.data : error.message);
      alert(`Failed to update profile:\n${error.response ? error.response.data.error : error.message}`);
    } finally {
      setCurrentPassword("");
      setNewPassword("");
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  const handleDelete = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }
      await axios.delete(`${apiBaseUrl}/api/v1/member`, {
        headers: {
          "X-CSRF-Token": csrfToken || "",
        },
        withCredentials: true,
      });
      navigate("/");
    } catch (error) {
      console.error(
        "Error deleting :",
        error.response ? error.response.data : error.message
      );
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
        location.reload();
      })
      .catch((error) => {
        console.log("Logout failed:", error);
      });
  };

  const handleSendMessage = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }

      // Send the message to the backend
      await axios.post(
        `${apiBaseUrl}/api/v1/messages`,
        {
          recipient: id,
          content: messageContent,
        },
        {
          headers: {
            "X-CSRF-Token": csrfToken,
          },
          withCredentials: true,
        }
      );

      alert("Message sent successfully!");
      setMessageContent("");
      setIsMessageBoxOpen(false);
    } catch (error) {
      console.error("Error sending message:", error.response ? error.response.data : error.message);
      alert(`Failed to send message:\n${error.response ? error.response.data.error : error.message}`);
    }
  };

  const [posts, setPosts] = useState([]);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await axios.get(
          `${apiBaseUrl}/api/v1/member/${id}/posts`
        );
        setPosts(response.data.data);
      } catch {
        setPosts([]);
      }
    };
    fetchPosts();
  }, [apiBaseUrl, id]);
  const { totalLikes, totalDislikes } = calcLikesAndDislikes(posts);

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        {/* Profile Card */}
        <div className="card bg-base-100 shadow-xl mb-6 w-full max-w-md">
          <div className="card-body pt-6 pb-0">
            {/* Card Header: Avatar, Name, Email, and Edit Button */}
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className="avatar w-16 h-16 rounded-full flex justify-center items-center bg-primary text-3xl mr-4">
                  {id.trim().charAt(0).toUpperCase()}
                </div>
                <div>
                  <h2 className="card-title">{id}</h2>
                  <p className="flex gap-2">
                    <div className="badge badge-primary text-xs">
                      <ThumbsUp className="w-[1em]" /> {totalLikes}
                    </div>
                    <div className="badge badge-primary text-xs">
                      <ThumbsDown className="w-[1em]" /> {totalDislikes}
                    </div>
                  </p>
                </div>
              </div>
            </div>
            {/* Card Body: Bio or Edit Form */}
            {isEditing ? (
              <div className="w-full flex flex-col gap-4 mt-2">
                <div className="form-control">
                  <input
                    id="username"
                    placeholder="Username"
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="input input-bordered w-full"
                  />
                </div>
                <div className="form-control">
                  <input
                    id="email"
                    type="email"
                    placeholder="example@ufl.edu"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="input input-bordered w-full"
                  />
                </div>
                <div className="form-control">
                  <input
                    id="currentPassword"
                    type="password"
                    placeholder="Current Password"
                    value={currentPassword}
                    onChange={(e) => setCurrentPassword(e.target.value)}
                    className="input input-bordered w-full"
                  />
                </div>
                <div className="form-control">
                  <input
                    id="newPassword"
                    type="password"
                    placeholder="New Password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    className="input input-bordered w-full"
                  />
                </div>
                <div className="form-control">
                  <textarea
                    id="bio"
                    value={bio}
                    placeholder="Tell us something about yourself..."
                    onChange={(e) => setBio(e.target.value)}
                    className="textarea textarea-bordered w-full"
                    rows="3"
                  ></textarea>
                </div>
                {/* Card Footer: Delete, Logout, Save, Cancel */}
                <div className="flex flex-wrap justify-end gap-2">
                  <button
                    onClick={handleDelete}
                    className="btn btn-error btn-sm flex-1"
                  >
                    Delete Profile
                  </button>
                  <button
                    onClick={handleLogout}
                    className="btn btn-secondary btn-sm flex-1"
                  >
                    Logout
                  </button>
                  <button
                    onClick={handleSave}
                    className="btn btn-primary btn-sm flex-1"
                  >
                    Save Changes
                  </button>
                </div>
              </div>
            ) : (
              <p className="mt-2">{bio}</p>
            )}
            {/* Show Message button if not own profile */}
            {ownProfile ? (
              // Only show the Edit button when not editing
              isEditing ? (
                <button
                  onClick={handleCancel}
                  className="btn btn-outline btn-sm mb-6"
                >
                  Cancel
                </button>
              ) : (
                <button
                  onClick={() => setIsEditing(true)}
                  className="btn btn-outline btn-sm mb-6"
                  title="Edit Profile"
                >
                  Edit
                </button>
              )
            ) : (
              <div></div>
            )}
          </div>
        </div>
        {isMessageBoxOpen && (
            <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
              <div className="bg-white text-black p-6 rounded shadow-lg w-full max-w-md z-50">
                <h3 className="text-lg font-bold mb-4">Send a Message to {id}</h3>
                <textarea
            value={messageContent}
            onChange={(e) => setMessageContent(e.target.value)}
            className="textarea textarea-bordered w-full mb-4"
            rows="4"
            placeholder="Type your message here..."
                ></textarea>
                <div className="flex justify-end gap-2">
            <button
              onClick={() => setIsMessageBoxOpen(false)}
              className="btn btn-secondary btn-sm"
            >
              Cancel
            </button>
            <button
              onClick={handleSendMessage}
              className="btn btn-primary btn-sm"
            >
              Send
            </button>
                </div>
              </div>
            </div>
          )}
          {/* Posts List */}
        <div className="w-full flex flex-col justify-center items-center">
          {posts.length > 0 ? (
            posts.map((post) => (
              <PostCard key={post.post_id} post={post} preview={true} />
            ))
          ) : (
            <p className="text-lg text-accent-content">No posts found</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Home;
