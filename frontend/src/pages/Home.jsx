import React from "react";
import { useParams } from "react-router-dom";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";

const User = () => {
  const { searchTerm } = useParams();
  console.log("Search term:", searchTerm);

  const posts = [
    {
      id: 31,
      title: "Post Title 1",
      author: "Author",
      date: "2/4/25",
      content:
        "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sed exercitationem asperiores ratione laudantium nihil eius pariatur, laboriosam aliquam, voluptatem a ab. Saepe maxime minus, fugit deserunt beatae nihil sint optio possimus reiciendis debitis sit error quam, tenetur doloribus. Eum iusto debitis suscipit dolores, id excepturi itaque ipsam eaque earum magni ea tempore, asperiores expedita at quibusdam repellendus, illum praesentium quas? Quae impedit in nesciunt ratione ut similique vero maiores itaque perspiciatis, odit accusantium dignissimos minima illo repellat quasi optio cumque ducimus. Minus blanditiis tempora nihil quidem a non. Obcaecati, aliquam nostrum voluptatibus, adipisci consequuntur architecto nisi, modi provident magnam fugiat in. Similique fuga quam quibusdam accusantium voluptatem sint delectus, quisquam nobis veritatis velit blanditiis, itaque tempora eum minus magnam voluptas est modi dolorem provident recusandae laudantium perspiciatis. Nam, harum labore sit sunt nobis repellendus voluptatum optio aspernatur voluptas, debitis dicta? Nesciunt exercitationem provident quisquam eaque consequuntur laborum. Modi, molestiae ipsum repudiandae reprehenderit at alias sunt.",
      images: [
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
        "https://images.pexels.com/photos/1054666/pexels-photo-1054666.jpeg",
        "https://img.daisyui.com/images/stock/photo-1559703248-dcaaec9fab78.webp",
      ],
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 2,
          author: "User1",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 23,
          author: "User2",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 2,
          author: "User3",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
      ],
    },
    {
      id: 2,
      title: "Post Title 2",
      author: "Author",
      date: "2/4/25",
      content:
        "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sed exercitationem asperiores ratione laudantium nihil eius pariatur, laboriosam aliquam, voluptatem a ab. Saepe maxime minus, fugit deserunt beatae nihil sint optio possimus reiciendis debitis sit error quam, tenetur doloribus. Eum iusto debitis suscipit dolores, id excepturi itaque ipsam eaque earum magni ea tempore, asperiores expedita at quibusdam repellendus, illum praesentium quas? Quae impedit in nesciunt ratione ut similique vero maiores itaque perspiciatis, odit accusantium dignissimos minima illo repellat quasi optio cumque ducimus. Minus blanditiis tempora nihil quidem a non. Obcaecati, aliquam nostrum voluptatibus, adipisci consequuntur architecto nisi, modi provident magnam fugiat in. Similique fuga quam quibusdam accusantium voluptatem sint delectus, quisquam nobis veritatis velit blanditiis, itaque tempora eum minus magnam voluptas est modi dolorem provident recusandae laudantium perspiciatis. Nam, harum labore sit sunt nobis repellendus voluptatum optio aspernatur voluptas, debitis dicta? Nesciunt exercitationem provident quisquam eaque consequuntur laborum. Modi, molestiae ipsum repudiandae reprehenderit at alias sunt.",
      images: [],
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 2,
          author: "User1",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 3,
          author: "User2",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 4,
          author: "User3",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
      ],
    },
  ];

  const filteredPosts = searchTerm
    ? posts.filter((post) => post.title.toLowerCase().includes(searchTerm.toLowerCase()))
    : posts;


  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        {filteredPosts.length > 0 ? (
          filteredPosts.map((post) => <PostCard key={post.id} post={post} preview={true} />)
        ) : (
          <p className="text-lg text-gray-500">No posts found</p>
        )}
      </div>
    </div>
  );
};

export default User;
