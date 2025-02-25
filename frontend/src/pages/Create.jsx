import React, { useState, useRef, useEffect } from "react";
import NavBar from "./components/NavBar";
import { Paperclip } from "lucide-react";

const Create = () => {
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");
  const [images, setImages] = useState([]);
  const imageContainerRef = useRef(null);
  const fileInputRef = useRef(null);

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

  const handleCreatePost = () => {
    alert(`Title:\n${title}\n\nText:\n${text}`);
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
                onClick={() => handleCreatePost()}
              >
                Create Post
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

export default Create;
