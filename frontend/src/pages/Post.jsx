import React, { useEffect, useState } from "react";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";
import CommentCard from "./components/CommentCard";
import { useParams } from "react-router-dom";

const Post = () => {
  const [post, setPost] = useState([]);
  const { id } = useParams();

  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await axios.get(
          `${apiBaseUrl}/api/v1/post/${id}`
        );
        setPost(response.data.data);
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    };
    fetchPosts();
  }, []);

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        <PostCard post={post} />
        {post.comments &&
          post.comments.map((comment) => (
            <CommentCard key={comment.comment_id} comment={comment} />
          ))}
      </div>
    </div>
  );
};

export default Post;
