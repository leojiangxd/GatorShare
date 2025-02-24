import React from "react";
import NavBar from "./components/NavBar";
import PostCard from "./components/PostCard";

const User = () => {
  const posts = [
    {
      id: 2,
      title: "Post Title",
      author: "Author",
      date: "2/4/25",
      content:
        "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sed exercitationem asperiores ratione laudantium nihil eius pariatur, laboriosam aliquam, voluptatem a ab. Saepe maxime minus, fugit deserunt beatae nihil sint optio possimus reiciendis debitis sit error quam, tenetur doloribus. Eum iusto debitis suscipit dolores, id excepturi itaque ipsam eaque earum magni ea tempore, asperiores expedita at quibusdam repellendus, illum praesentium quas? Quae impedit in nesciunt ratione ut similique vero maiores itaque perspiciatis, odit accusantium dignissimos minima illo repellat quasi optio cumque ducimus. Minus blanditiis tempora nihil quidem a non. Obcaecati, aliquam nostrum voluptatibus, adipisci consequuntur architecto nisi, modi provident magnam fugiat in. Similique fuga quam quibusdam accusantium voluptatem sint delectus, quisquam nobis veritatis velit blanditiis, itaque tempora eum minus magnam voluptas est modi dolorem provident recusandae laudantium perspiciatis. Nam, harum labore sit sunt nobis repellendus voluptatum optio aspernatur voluptas, debitis dicta? Nesciunt exercitationem provident quisquam eaque consequuntur laborum. Modi, molestiae ipsum repudiandae reprehenderit at alias sunt.",
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 1,
          author: "User1",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 1,
          author: "User2",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 1,
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
      title: "Post Title",
      author: "Author2",
      date: "2/4/25",
      content:
        "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Sed exercitationem asperiores ratione laudantium nihil eius pariatur, laboriosam aliquam, voluptatem a ab. Saepe maxime minus, fugit deserunt beatae nihil sint optio possimus reiciendis debitis sit error quam, tenetur doloribus. Eum iusto debitis suscipit dolores, id excepturi itaque ipsam eaque earum magni ea tempore, asperiores expedita at quibusdam repellendus, illum praesentium quas? Quae impedit in nesciunt ratione ut similique vero maiores itaque perspiciatis, odit accusantium dignissimos minima illo repellat quasi optio cumque ducimus. Minus blanditiis tempora nihil quidem a non. Obcaecati, aliquam nostrum voluptatibus, adipisci consequuntur architecto nisi, modi provident magnam fugiat in. Similique fuga quam quibusdam accusantium voluptatem sint delectus, quisquam nobis veritatis velit blanditiis, itaque tempora eum minus magnam voluptas est modi dolorem provident recusandae laudantium perspiciatis. Nam, harum labore sit sunt nobis repellendus voluptatum optio aspernatur voluptas, debitis dicta? Nesciunt exercitationem provident quisquam eaque consequuntur laborum. Modi, molestiae ipsum repudiandae reprehenderit at alias sunt.",
      likes: 5125,
      dislikes: 5125,
      views: 10521,
      comments: [
        {
          id: 1,
          author: "User1",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 1,
          author: "User2",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
        {
          id: 1,
          author: "User3",
          date: "2/4/25",
          content: "Nice post!",
          likes: 10,
          dislikes: 10,
        },
      ],
    },
  ];

  return (
    <div className="flex flex-col h-screen bg-base-200">
      <NavBar />
      <div className="p-10 overflow-y-auto flex-grow flex flex-col items-center">
        {posts.map((post) => (
          <PostCard key={post.id} post={post} preview={true} />
        ))}
      </div>
    </div>
  );
};

export default User;
