import React from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import {Home, Login, Post, Register, User, Create, Notifications} from './pages/pages.js'

const App = () => {
  return (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/search/:searchTerm" element={<Home />} />
			<Route path="/login" element={<Login />} />
			<Route path="/register" element={<Register />} />
			<Route path="/post/:id" element={<Post />} />
			<Route path="/user/:id" element={<User />} />
			<Route path="/create/" element={<Create />} />
			<Route path="/notifications" element={<Notifications />} />

			{/* Navigate to / if parameters don't exist */}
			<Route path="/search" element={<Navigate to="/" />} />
			<Route path="/post/" element={<Navigate to="/" />} />
			<Route path="/user/" element={<Navigate to="/" />} />
		</Routes>
	);
};

export default App;