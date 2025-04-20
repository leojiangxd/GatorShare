import React, { useEffect, useState } from "react";
import { ChevronDown, ThumbsDown, ThumbsUp } from "lucide-react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { formatTime, getCsrfToken, getUsername } from "../../utils/functions";
import axios from "axios";

const CommentCard = ({ comment }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editedContent, setEditedContent] = useState(comment.content);
  const { id } = useParams();
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  const navigate = useNavigate();
  const [ownComment, setOwnComment] = useState(false);
  
  // New state for like/dislike functionality
  const [likes, setLikes] = useState(comment.likes || 0);
  const [dislikes, setDislikes] = useState(comment.dislikes || 0);
  const [isLiked, setIsLiked] = useState(false);
  const [isDisliked, setIsDisliked] = useState(false);

  useEffect(() => {
    const initializeComment = async () => {
      try {
        const username = await getUsername();
        if (!username) return;
        
        if (username === comment.author) {
          setOwnComment(true);
        }
        
        // Initialize likes/dislikes count
        setLikes(comment.likes || 0);
        setDislikes(comment.dislikes || 0);

        const csrfToken = getCsrfToken();

        // Fetch user's like/dislike status for this comment
        const [likeResponse, dislikeResponse] = await Promise.all([
          axios.get(`${apiBaseUrl}/api/v1/member/${username}/liked-comments`, {
            headers: {
              "X-CSRF-Token": csrfToken,
            },
            withCredentials: true,
          }),
          axios.get(`${apiBaseUrl}/api/v1/member/${username}/disliked-comments`, {
            headers: {
              "X-CSRF-Token": csrfToken,
            },
            withCredentials: true,
          }),
        ]);

        setIsLiked(
          likeResponse.data.data.some(
            (likedComment) => likedComment.comment_id === comment.comment_id
          )
        );

        setIsDisliked(
          dislikeResponse.data.data.some(
            (dislikedComment) => dislikedComment.comment_id === comment.comment_id
          )
        );
      } catch (error) {
        console.error("Error initializing comment:", error);
      }
    };
  
    initializeComment();
  }, [comment]); 

  const handleLike = async (e) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }

      const response = await axios.put(
        `${apiBaseUrl}/api/v1/comment/${id}/${comment.comment_id}/like-dislike`,
        { action: "like" },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );

      // Update the like/dislike counts from the response
      setLikes(response.data.likes);
      setDislikes(response.data.dislikes);
      setIsLiked(!isLiked);
      setIsDisliked(false);
    } catch (error) {
      console.error(
        "Error liking comment:",
        error.response ? error.response.data : error.message
      );
    }
  };

  const handleDislike = async (e) => {
    e.preventDefault();
    e.stopPropagation();

    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }

      const response = await axios.put(
        `${apiBaseUrl}/api/v1/comment/${id}/${comment.comment_id}/like-dislike`,
        { action: "dislike" },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );

      // Update the like/dislike counts from the response
      setLikes(response.data.likes);
      setDislikes(response.data.dislikes);
      setIsDisliked(!isDisliked);
      setIsLiked(false);
    } catch (error) {
      console.error(
        "Error disliking comment:",
        error.response ? error.response.data : error.message
      );
    }
  };

  const handleDelete = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }
      await axios.delete(
        `${apiBaseUrl}/api/v1/comment/${id}/${comment.comment_id}`,
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );
      window.location.reload();
    } catch (error) {
      console.error(
        "Error deleting :",
        error.response ? error.response.data : error.message
      );
    }
  };

  const handleEditSubmit = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }

      if (editedContent.trim() === "") {
        handleDelete();
        return;
      }

      const response = await axios.put(
        `${apiBaseUrl}/api/v1/comment/${id}/${comment.comment_id}`,
        { content: editedContent.trim() },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );
      setIsEditing(false);
      window.location.reload();
    } catch (error) {
      console.error(
        "Error updating comment:",
        error.response ? error.response.data : error.message
      );
    }
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
              <span className="opacity-50 ml-1">
                {formatTime(comment.CreatedAt)}
              </span>
              {ownComment ? (
                <div className="dropdown">
                  <div
                    tabIndex={0}
                    role="button"
                    className="rounded-full cursor-pointer"
                  >
                    <ChevronDown className="h-[1em]" />
                  </div>
                  <ul
                    tabIndex={0}
                    className="dropdown-content menu bg-accent-content rounded-box z-1 w-52 p-2 shadow-sm"
                  >
                    <li>
                      <a onClick={() => setIsEditing(!isEditing)}>
                        {isEditing ? "Cancel" : "Edit"}
                      </a>
                    </li>
                    <li>
                      <a onClick={handleDelete}>Delete</a>
                    </li>
                  </ul>
                </div>
              ) : (
                <></>
              )}
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleLike}
                className={
                  isLiked
                    ? "liked badge badge-secondary flex items-center cursor-pointer text-xs"
                    : "badge badge-primary hover:badge-secondary flex items-center cursor-pointer transition-colors duration-350 text-xs"
                }
              >
                <ThumbsUp className="w-[1em]" /> {likes}
              </button>
              <button
                onClick={handleDislike}
                className={
                  isDisliked
                    ? "disliked badge badge-secondary flex items-center cursor-pointer text-xs"
                    : "badge badge-primary hover:badge-secondary flex items-center cursor-pointer transition-colors duration-350 text-xs"
                }
              >
                <ThumbsDown className="w-[1em]" /> {dislikes}
              </button>
            </div>
          </div>
        </div>

        {isEditing ? (
          <div className="flex flex-col gap-2">
            <input
              type="text"
              placeholder="Comment"
              className="input w-full"
              value={editedContent}
              onChange={(e) => setEditedContent(e.target.value)}
            />
            <div className="flex flex-row justify-end gap-2">
              <button
                onClick={() => setIsEditing(!isEditing)}
                className="btn btn-sm"
              >
                Cancel
              </button>
              <button
                onClick={handleEditSubmit}
                className="btn btn-primary btn-sm"
              >
                Save
              </button>
            </div>
          </div>
        ) : (
          <div>{comment.content}</div>
        )}
      </div>
    </div>
  );
};

export default CommentCard;