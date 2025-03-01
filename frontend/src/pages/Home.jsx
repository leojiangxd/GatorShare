import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";

const User = () => {
  const { searchTerm } = useParams();

  const [posts, setPosts] = useState([]);
  
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await axios.get(`${apiBaseUrl}/api/v1/post`);
        setPosts(response.data.data);
      } catch {
        setPosts([])
      }
    };
    fetchPosts();
  }, []);

  const filteredPosts = searchTerm
    ? posts.filter((post) => post.title.toLowerCase().includes(searchTerm.toLowerCase()))
    : posts;

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        {filteredPosts.length > 0 ? (
          filteredPosts.map((post) => <PostCard key={post.post_id} post={post} preview={true} />)
        ) : (
          <p className="text-lg text-gray-500">No posts found</p>
        )}
      </div>
    </div>
  );
};

export default User;
