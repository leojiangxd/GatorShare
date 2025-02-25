import React, { useState } from "react";
import { Link, useParams } from "react-router-dom";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";
import { ThumbsDown, ThumbsUp } from "lucide-react";

const User = () => {
  const loggedInUsername = "test";
  const { id } = useParams();
  const ownProfile =
    id.toLowerCase().trim() == loggedInUsername.toLowerCase().trim();

  // Toggle between view and edit modes.
  const [isEditing, setIsEditing] = useState(false);

  // Profile state.
  const [username, setUsername] = useState(id);
  const [email, setEmail] = useState("example@example.com");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [bio, setBio] = useState("Tell us something about yourself...");

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
      id: 4,
      title: "Post Title",
      author: "Author",
      date: "2/4/25",
      content: "Lorem ipsum dolor sit amet consectetur, adipisicing elit...",
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 5,
          author: "User1",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        // ...other comments
      ],
    },
    {
      id: 10,
      title: "Post Title",
      author: "Author",
      date: "2/4/25",
      content: "Lorem ipsum dolor sit amet consectetur, adipisicing elit...",
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 6,
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
  const { totalLikes, totalDislikes } = calcLikesAndDislikes(posts);

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        {/* Profile Card */}
        <div className="card bg-base-100 shadow-xl mb-6 w-full max-w-md">
          <div className="card-body pt-6">
            {/* Card Header: Avatar, Name, Email, and Edit Button */}
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className="avatar mr-4">
                  <div className="w-16 rounded-full">
                    <img
                      src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"
                    />
                  </div>
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
                  className="btn btn-outline btn-sm"
                >
                  Cancel
                </button>
              ) : (
                <button
                  onClick={() => setIsEditing(true)}
                  className="btn btn-outline btn-sm"
                  title="Edit Profile"
                >
                  Edit
                </button>
              )
            ) : (
              <Link to={`/message/${id}`} className="btn btn-primary btn-sm">
                Message
              </Link>
            )}
          </div>
        </div>
        {/* Posts List */}
        <div className="w-full flex flex-col justify-center items-center">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} preview={true} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default User;
