import { Eye, Heart, MessageSquare, ThumbsDown, ThumbsUp } from "lucide-react";
import { Link } from "react-router-dom";

const PostCard = ({ post, preview = false }) => {
  const handleLike = (e) => {
    e.preventDefault();
		alert('liked post');
  };
  const handleDislike = (e) => {
    e.preventDefault();
		alert('disliked post');
  };

  const cardContent = (
    <>
      <h2 className="card-title">{post.title}</h2>
      <span className="text-xs flex">
        <div className="flex flex-auto items-center">
          <Link to={`/user/${post.author}`} className="hover:underline">
            {post.author}
          </Link>
          <span className="opacity-50 ml-1">{post.date}</span>
        </div>
        <div className="card-actions flex items-center">
          <div className="badge badge-ghost flex items-center text-xs">
            <MessageSquare className="w-[1em]" /> {post.comments.length}
          </div>
          <div className="badge badge-ghost flex items-center text-xs">
            <Eye className="w-[1em]" /> {post.views}
          </div>
          <button
            onClick={handleLike}
            className="badge badge-primary hover:badge-secondary flex items-center cursor-pointer transition-colors duration-350 text-xs"
          >
            <ThumbsUp className="w-[1em]" /> {post.likes}
          </button>
          <button
            onClick={handleDislike}
            className="badge badge-primary hover:badge-secondary flex items-center cursor-pointer transition-colors duration-350 text-xs"
          >
            <ThumbsDown className="w-[1em]" /> {post.dislikes}
          </button>
        </div>
      </span>
      {preview ? (
        <p>
          {post.content.split(" ").slice(0, 50).join(" ")}
          {post.content.split(" ").length > 50 && (
            <div className="text-center">. . .</div>
          )}
        </p>
      ) : (
        <p>{post.content}</p>
      )}
    </>
  );

  return (
    <div className="card w-full bg-base-100 card-md shadow-sm mb-5 max-w-300">
      {preview ? (
        <Link className="card-body" to={`/post/${post.id}`}>
          {cardContent}
        </Link>
      ) : (
        <div className="card-body">{cardContent}</div>
      )}
    </div>
  );
};

export default PostCard;
