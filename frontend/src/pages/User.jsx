import React, { useState } from "react";
import { useParams } from "react-router-dom";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";

const User = () => {
  const { id } = useParams();

  // Toggle between view and edit modes.
  const [isEditing, setIsEditing] = useState(false);

  // Profile state.
  const [username, setUsername] = useState(id || "User");
  const [email, setEmail] = useState("example@example.com");
  const [bio, setBio] = useState("Tell us something about yourself...");

  // Handlers to simulate actions.
  const handleSave = () => {
    alert("Profile updated!");
    setIsEditing(false);
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  const handleDelete = () => {
    alert("Profile deleted!");
  };

  const handleLogout = () => {
    alert("Logged out!");
  };

  // Dummy posts (replace or fetch as needed)
  const posts = [
    {
      id: 1,
      title: "Post Title",
      author: "Author",
      date: "2/4/25",
      content: "Lorem ipsum dolor sit amet consectetur, adipisicing elit...",
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 1,
          author: "User1",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        // ...other comments
      ],
    },
    // ...more posts
  ];

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        {/* Profile Card */}
        <div className="card bg-base-100 shadow-xl p-6 mb-6 w-full max-w-md">
          {/* Card Header: Avatar, Name, Email, and Edit Button */}
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <div className="avatar mr-4">
                <div className="w-16 rounded-full">
                  <img
                    alt="User avatar"
                    src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"
                  />
                </div>
              </div>
              <div>
                <h2 className="card-title">{username}</h2>
                <p className="text-sm text-gray-600">{email}</p>
              </div>
            </div>
            {/* Only show the Edit button when not editing */}
            {!isEditing && (
              <button
                onClick={() => setIsEditing(true)}
                className="btn btn-outline btn-sm"
                title="Edit Profile"
              >
                Edit
              </button>
            )}
          </div>

          {/* Card Body: Bio or Edit Form */}
          <div className="card-body pt-4">
            {isEditing ? (
              <div className="w-full">
                <div className="form-control mb-4">
                  <label className="label" htmlFor="username">
                    <span className="label-text">Username</span>
                  </label>
                  <input
                    id="username"
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="input input-bordered w-full"
                  />
                </div>
                <div className="form-control mb-4">
                  <label className="label" htmlFor="email">
                    <span className="label-text">Email</span>
                  </label>
                  <input
                    id="email"
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="input input-bordered w-full"
                  />
                </div>
                <div className="form-control mb-4">
                  <label className="label" htmlFor="bio">
                    <span className="label-text">Bio</span>
                  </label>
                  <textarea
                    id="bio"
                    value={bio}
                    onChange={(e) => setBio(e.target.value)}
                    className="textarea textarea-bordered w-full"
                    rows="3"
                  ></textarea>
                </div>
                <div className="flex flex-wrap gap-2">
                  <button
                    onClick={handleSave}
                    className="btn btn-primary btn-sm"
                  >
                    Save Changes
                  </button>
                  <button
                    onClick={handleCancel}
                    className="btn btn-outline btn-sm"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            ) : (
              <p className="mt-2">{bio}</p>
            )}

            {/* Card Footer: Delete and Logout */}
            <div className="flex justify-end mt-4 gap-2">
              <button
                onClick={handleDelete}
                className="btn btn-error btn-sm"
              >
                Delete Profile
              </button>
              <button
                onClick={handleLogout}
                className="btn btn-secondary btn-sm"
              >
                Logout
              </button>
            </div>
          </div>
        </div>

        {/* Posts List */}
        <div className="w-full">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} preview={true} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default User;
