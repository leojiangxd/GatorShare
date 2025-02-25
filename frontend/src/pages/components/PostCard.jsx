import { useEffect, useRef, useState } from "react";
import { Eye, MessageSquare, Send, ThumbsDown, ThumbsUp } from "lucide-react";
import { Link } from "react-router-dom";

const PostCard = ({ post, preview = false }) => {
  const [comment, setComment] = useState("");
  const imageContainerRef = useRef(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [selectedImage, setSelectedImage] = useState("");

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

  const handleLike = (e) => {
    e.preventDefault();
    e.stopPropagation();
    alert("liked post");
  };

  const handleDislike = (e) => {
    e.preventDefault();
    e.stopPropagation();
    alert("disliked post");
  };

  const handleSendComment = () => {
    alert(`Comment: ${comment}`);
    setComment("");
  };

  const incrementViewCount = () => {
    console.log(`Views: ${post.views + 1}`);
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
      <h2 className="card-title">{post.title}</h2>
      <span className="text-xs flex flex-wrap">
        <div className="flex flex-auto items-center">
          <Link to={`/user/${post.author}`} className="hover:underline">
            {post.author}
          </Link>
          <span className="opacity-50 ml-1">{post.date}</span>
        </div>
        <div className="card-actions flex items-center pt-2">
          <div className="badge badge-ghost flex text-xs">
            <MessageSquare className="w-[1em]" /> {post.comments.length}
          </div>
          <div className="badge badge-ghost flex text-xs">
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
      <div className="flex justify-center">
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
        <div className="flex justify-center items-center w-auto mt-5">
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
            to={`/post/${post.id}`}
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
