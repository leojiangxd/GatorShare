import React, { useState, useRef, useEffect } from "react";
import NavBar from "./components/NavBar";
import { Paperclip } from "lucide-react";
import { getCsrfToken, getUsername } from "../utils/functions";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

const EditPost = () => {
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");
  const [images, setImages] = useState([]);
  const imageContainerRef = useRef(null);
  const fileInputRef = useRef(null);

  const navigate = useNavigate();
  const { id } = useParams();

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const currentUsername = await getUsername();
        const response = await axios.get(
          `${apiBaseUrl}/api/v1/post/${id}`
        );

        const { title, content, images: postImages, author: postAuthor } = response.data.data;

        // Redirect if current user is not the author
        if (currentUsername !== postAuthor) {
          console.log(currentUsername)
          // navigate(`/post/${id}`);
          return;
        }
        
        setTitle(title);
        setText(content);
        setImages(postImages || []);
      } catch (error) {
        console.error("Error fetching post:", error);
        navigate(`/post/${id}`);
      }
    };

    fetchPost();
  }, [id, navigate]);

  const handleImageChange = (e) => {
    const files = Array.from(e.target.files);
    const validImages = files.filter((file) => file.type.startsWith("image/"));
    const imageUrls = validImages.map((file) => {
      return new Promise((resolve) => {
        const reader = new FileReader();
        reader.onloadend = () => resolve(reader.result);
        reader.readAsDataURL(file);
      });
    });
    Promise.all(imageUrls).then((results) =>
      setImages((images) => [...images, ...results])
    );
  };

  const handleImageClick = (index) => {
    setImages((images) => images.filter((_, i) => i !== index));
  };

  const handleUpdatePost = async () => {
    if (title.trim() == "" || text.trim() == "") return;
    const csrfToken = getCsrfToken();
    
    if (!csrfToken) {
      console.error("CSRF token is missing.");
      navigate("/login");
      return;
    }
  
    try {
      await axios.put(
        `${apiBaseUrl}/api/v1/post/${id}`,
        { title: title.trim(), content: text.trim(), images },
        {
          headers: {
            "X-CSRF-Token": csrfToken || "",
          },
          withCredentials: true,
        }
      );
      navigate(`/post/${id}`);
    } catch (error) {
      console.error("Error updating post:", error.response ? error.response.data : error.message);
    }
  };

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
  }, [images]);

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="flex flex-1 items-center justify-center m-10">
        <div className="card h-full w-full bg-base-100 shadow-sm max-w-300">
          <div className="card-body flex flex-col gap-5">
            <div className="flex gap-5">
              <input
                type="text"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="input flex-grow"
              />
              <button
                className="btn btn-secondary"
                onClick={() => fileInputRef.current.click()}
              >
                <Paperclip className="w-[1em]" />
              </button>
              <button
                className="btn btn-primary"
                onClick={() => handleUpdatePost()}
              >
                Update Post
              </button>
            </div>
            <textarea
              placeholder="Type here"
              value={text}
              onChange={(e) => setText(e.target.value)}
              className="textarea w-full flex-grow"
            ></textarea>
            <input
              type="file"
              accept="image/*"
              multiple
              onChange={handleImageChange}
              className="hidden"
              ref={fileInputRef}
            />
            {/* Image container */}
            {images.length > 0 && (
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
                        onClick={() => handleImageClick(index)}
                      />
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default EditPost;
