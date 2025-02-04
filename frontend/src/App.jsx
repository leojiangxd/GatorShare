import React from "react";
import { Route, Routes } from "react-router-dom";
import {Home, Login, Register} from './pages/pages.js'

const App = () => {
  return (
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/login" element={<Login />} />
			<Route path="/register" element={<Register />} />
		</Routes>
	);
};

export default App;