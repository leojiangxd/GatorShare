import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";
import axios from "axios";
import { ArrowLeft, ArrowRight, ChevronDown } from "lucide-react";

const User = () => {
  const { searchTerm } = useParams();
  const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

  const [posts, setPosts] = useState([]);
  const [column, setColumn] = useState("created_at");
  const [order, setOrder] = useState("DESC");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);

  const fetchPosts = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${apiBaseUrl}/api/v1/post`, {
        params: {
          column,
          order,
          limit: 10,
          offset: (page - 1) * 10,
          search_key: searchTerm || "",
        },
      });
      setPosts(response.data.data);
      setTotalPages(Math.ceil(response.data.count / limit));
    } catch (error) {
      console.error("Error fetching posts:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, [column, order, searchTerm, page]);

  const handleNextPage = () => {
    if (page < totalPages) {
      setPage((prevPage) => prevPage + 1);
    }
  };

  const handlePreviousPage = () => {
    if (page > 1) {
      setPage((prevPage) => prevPage - 1);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        <div className="flex gap-1 mb-10">
          <details className="dropdown">
            <summary className="btn bg-accent-content m-1">
              {column === "created_at"
                ? "Created At"
                : column.charAt(0).toUpperCase() + column.slice(1)}
              <ChevronDown className="w-[1em]" />
            </summary>
            <ul className="menu dropdown-content bg-accent-content rounded-box z-1 w-52 p-2 shadow-sm">
              <li>
                <a onClick={() => setColumn("created_at")}>Created At</a>
              </li>
              <li>
                <a onClick={() => setColumn("likes")}>Likes</a>
              </li>
              <li>
                <a onClick={() => setColumn("dislikes")}>Dislikes</a>
              </li>
              <li>
                <a onClick={() => setColumn("views")}>Views</a>
              </li>
              <li>
                <a onClick={() => setColumn("comments")}>Comments</a>
              </li>
            </ul>
          </details>
          <details className="dropdown">
            <summary className="btn bg-accent-content m-1">
              {order === "DESC" ? "Descending" : "Ascending"}
              <ChevronDown className="w-[1em]" />
            </summary>
            <ul className="menu dropdown-content bg-accent-content rounded-box z-1 w-52 p-2 shadow-sm">
              <li>
                <a onClick={() => setOrder("ASC")}>Ascending</a>
              </li>
              <li>
                <a onClick={() => setOrder("DESC")}>Descending</a>
              </li>
            </ul>
          </details>
        </div>

        {loading ? (
          <></>
        ) : posts.length > 0 ? (
          <>
            {posts.map((post) => (
              <PostCard key={post.post_id} post={post} preview={true} />
            ))}
            <div className="flex justify-center gap-4 mt-4">
              <button
                className="btn btn-sm bg-primary"
                onClick={handlePreviousPage}
                disabled={page === 1}
              >
                <ArrowLeft className="w-[1em]" />
              </button>
              {page}/{totalPages}
              <button
                className="btn btn-sm bg-primary"
                onClick={handleNextPage}
                disabled={page === totalPages}
              >
                <ArrowRight className="w-[1em]" />
              </button>
            </div>
          </>
        ) : (
          <p className="text-lg text-accent-content">No posts found</p>
        )}
      </div>
    </div>
  );
};

export default User;
