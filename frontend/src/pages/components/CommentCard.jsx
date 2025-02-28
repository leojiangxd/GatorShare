import { ThumbsDown, ThumbsUp } from "lucide-react";
import { Link } from "react-router-dom";
import { formatTime } from "../../utils/functions";

const CommentCard = ({ comment }) => {
  const handleLike = (e) => {
    e.preventDefault();
    alert("liked comment");
  };
  const handleDislike = (e) => {
    e.preventDefault();
    alert("disliked comment");
  };

  return (
    <div className="card w-full bg-base-100 card-md shadow-sm mb-5 max-w-300">
      <div className="card-body">
        <div className="flex">
          <div className="chat-header flex items-center w-full flex-wrap">
            <div className="flex flex-auto items-center">
              <Link to={`/user/${comment.author}`} className="hover:underline">
                {comment.author}
              </Link>
              <span className="opacity-50 ml-1">{formatTime(comment.CreatedAt)}</span>
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleLike}
                className="badge badge-primary hover:badge-secondary flex items-center cursor-pointer transition-colors duration-350 text-xs"
              >
                <ThumbsUp className="w-[1em]" /> {comment.likes}
              </button>
              <button
                onClick={handleDislike}
                className="badge badge-primary hover:badge-secondary flex items-center cursor-pointer transition-colors duration-350 text-xs"
              >
                <ThumbsDown className="w-[1em]" /> {comment.dislikes}
              </button>
            </div>
          </div>
        </div>
        <div>{comment.content}</div>
      </div>
    </div>
  );
};

export default CommentCard;
