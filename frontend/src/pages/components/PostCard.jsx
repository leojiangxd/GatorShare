import { useEffect, useRef, useState } from "react";
import {
  ChevronDown,
  Eye,
  MessageSquare,
  Send,
  ThumbsDown,
  ThumbsUp,
} from "lucide-react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { formatTime, getCsrfToken, getUsername } from "../../utils/functions";
import axios from "axios";

const PostCard = ({ post, preview = false }) => {
  const { id } = useParams();
  const [comment, setComment] = useState("");
  const imageContainerRef = useRef(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [selectedImage, setSelectedImage] = useState("");
  const [likes, setLikes] = useState(post.likes || 0);
  const [dislikes, setDislikes] = useState(post.dislikes || 0);
  const [isLiked, setIsLiked] = useState(false);
  const [isDisliked, setIsDisliked] = useState(false);
  const [ownPost, setOwnPost] = useState(false);

  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;
  const navigate = useNavigate();

  useEffect(() => {
    if (!post?.post_id) return;

    const initalizePost = async () => {
      setLikes(post.likes || 0);
      setDislikes(post.dislikes || 0);

      try {
        const username = await getUsername();
        if (!username) return;

        if (username === post.author) {
          setOwnPost(true);
        }

        const [likeResponse, dislikeResponse] = await Promise.all([
          axios.get(`${apiBaseUrl}/api/v1/member/${username}/liked-posts`),
          axios.get(`${apiBaseUrl}/api/v1/member/${username}/disliked-posts`),
        ]);

        setIsLiked(
          likeResponse.data.data.some(
            (likedPost) => likedPost.post_id === post.post_id
          )
        );

        setIsDisliked(
          dislikeResponse.data.data.some(
            (dislikedPost) => dislikedPost.post_id === post.post_id
          )
        );
      } catch (error) {
        console.error("Error fetching like/dislike:", error);
      }
    };

    initalizePost();
  }, [post]);

  // Scroll horizontally using vertical scroll
  useEffect(() => {
    const container = imageContainerRef.current;
    if (container) {
      const handleWheel = (event) => {
        event.preventDefault();
        container.scrollBy({
          left: event.deltaY * 1.5,
          behavior: "smooth",
        });
      };

      container.addEventListener("wheel", handleWheel, { passive: false });

      return () => {
        container.removeEventListener("wheel", handleWheel);
      };
    }
  }, []);

  // Handle escape key to close modal
  useEffect(() => {
    const handleEscKey = (event) => {
      if (event.key === "Escape" && modalOpen) {
        setModalOpen(false);
      }
    };

    document.addEventListener("keydown", handleEscKey);

    // Prevent body scrolling when modal is open
    if (modalOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }

    return () => {
      document.removeEventListener("keydown", handleEscKey);
      document.body.style.overflow = "";
    };
  }, [modalOpen]);

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

      const postId = preview ? post.post_id : id;

      const response = await axios.put(
        `${apiBaseUrl}/api/v1/post/${postId}/like-dislike`,
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
        "Error liking post:",
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

      const postId = preview ? post.post_id : id;

      const response = await axios.put(
        `${apiBaseUrl}/api/v1/post/${postId}/like-dislike`,
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
        "Error disliking post:",
        error.response ? error.response.data : error.message
      );
    }
  };

  const navigateEdit = () => {
    navigate(`/edit/post/${id}`);
  };

  const handleDelete = async () => {
    try {
      const csrfToken = getCsrfToken();
      if (!csrfToken) {
        console.error("CSRF token is missing.");
        navigate("/login");
        return;
      }
      await axios.delete(`${apiBaseUrl}/api/v1/post/${id}`, {
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

  const handleSendComment = async () => {
    if (comment.trim() == "") return;
    const csrfToken = getCsrfToken();

    if (!csrfToken) {
      console.error("CSRF token is missing.");
      navigate("/login");
      return;
    }

    try {
      console.log(comment);
      await axios.post(
        `${apiBaseUrl}/api/v1/comment/${post.post_id}`,
        { content: comment.trim() },
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
        "Error sending comment:",
        error.response ? error.response.data : error.message
      );
    }
  };

  const incrementViewCount = async () => {
    await axios.put(
      `${apiBaseUrl}/api/v1/post/${post.post_id}/increment-views`
    );
  };

  const handleImageClick = (e, image) => {
    e.preventDefault();
    e.stopPropagation();
    setSelectedImage(image);
    setModalOpen(true);
  };

  const closeModal = (e) => {
    if (e) {
      e.preventDefault();
      e.stopPropagation();
    }
    setModalOpen(false);
  };

  const images = post.images ? post.images : [];

  const cardContent = (
    <>
      <span className="text-xs flex flex-wrap">
        <div className="flex flex-auto items-center">
          <Link to={`/user/${post.author}`} className="hover:underline">
            {post.author}
          </Link>
          <span className="opacity-50 ml-1">{formatTime(post.CreatedAt)}</span>
          {!preview && ownPost ? (
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
                  <a onClick={navigateEdit}>Edit</a>
                </li>
                <li>
                  <a onClick={handleDelete}>Delete</a>
                </li>
              </ul>
            </div>
          ) : null}
        </div>
        <div className="card-actions flex items-center pt-2">
          <div className="badge badge-ghost flex text-xs">
            <MessageSquare className="w-[1em]" />{" "}
            {post.comments ? post.comments.length : 0}
          </div>
          <div className="badge badge-ghost flex text-xs">
            <Eye className="w-[1em]" /> {post.views}
          </div>
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
      </span>
      <h2 className="card-title">{post.title}</h2>
      {preview ? (
        <p>
          {post.content.slice(0, 500)}
          {post.content.length > 500 && (
            <span className="text-center block">&bull; &bull; &bull;</span>
          )}
        </p>
      ) : (
        <>
          <p>{post.content}</p>
        </>
      )}

      {/* Image container */}
      <div className={`flex justify-center ${images.length > 0 ? "mt-5" : ""}`}>
        <div
          id="imageContainer"
          ref={imageContainerRef}
          className="rounded-lg overflow-hidden inline-block"
        >
          <div className="flex space-x-2 max-h-64">
            {images.map((image, index) => (
              <img
                key={index}
                src={image}
                className="object-cover cursor-pointer h-64"
                onClick={(e) => handleImageClick(e, image)}
              />
            ))}
          </div>
        </div>
      </div>

      {!preview && (
        <div className="flex justify-center items-center w-auto">
          <input
            type="text"
            placeholder="Comment"
            className="input mr-5 flex-auto"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
          />
          <button
            onClick={handleSendComment}
            className="btn btn-circle btn-primary"
          >
            <Send className="w-[1em]" />
          </button>
        </div>
      )}
    </>
  );

  return (
    <>
      <div className="card w-full bg-base-100 card-md shadow-sm mb-5 max-w-300">
        {preview ? (
          <Link
            className="card-body"
            to={`/post/${post.post_id}`}
            onClick={incrementViewCount}
          >
            {cardContent}
          </Link>
        ) : (
          <div className="card-body">{cardContent}</div>
        )}
      </div>

      {modalOpen && (
        <div
          className="fixed inset-0 flex items-center justify-center z-100"
          style={{ backgroundColor: "rgba(0, 0, 0, 0.5)" }}
          onClick={closeModal}
        >
          <div className="relative max-w-full max-h-full flex justify-center items-center">
            <img
              src={selectedImage}
              className="max-w-full max-h-screen object-contain p-5"
            />
          </div>
        </div>
      )}
    </>
  );
};

export default PostCard;
