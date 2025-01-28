import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import {HomePage, LoginPage, RegisterPage} from './pages/pages.jsx'

const App = () => {
  return (
		<Routes>
			<Route path="/" element={<HomePage />} />
			<Route path="/login" element={<LoginPage />} />
			<Route path="/register" element={<RegisterPage />} />
		</Routes>
	);
};

export default App;