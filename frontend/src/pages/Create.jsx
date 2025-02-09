import React, { useState } from "react";
import NavBar from "./components/NavBar";

const Create = () => {
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");

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
                className="btn btn-primary"
                onClick={() => alert(`Title:\n${title}\n\nText:\n${text}`)}
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
          </div>
        </div>
      </div>
    </div>
  );
};

export default Create;
