import React from "react";
import { Route, Routes } from "react-router-dom";
import {Home, Login, Post, Register, User} from './pages/pages.js'

const App = () => {
  return (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/login" element={<Login />} />
			<Route path="/register" element={<Register />} />
			<Route path="/post/:id" element={<Post />} />
			<Route path="/user/:id" element={<User />} />
		</Routes>
	);
};

export default App;