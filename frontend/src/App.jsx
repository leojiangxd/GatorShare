import React from "react";
import { Route, Routes } from "react-router-dom";
import {Home, Login, Post, Register, User, Create} from './pages/pages.js'

const App = () => {
  return (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/login" element={<Login />} />
			<Route path="/register" element={<Register />} />
			<Route path="/post/:id" element={<Post />} />
			<Route path="/user/:id" element={<User />} />
			<Route path="/create/" element={<Create />} />
		</Routes>
	);
};

export default App;