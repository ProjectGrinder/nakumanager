"use client";

import Link from "next/link";
import { users } from "../../Database";
import { useState } from "react";
export default function Register() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");

  const handleRegister = () => {
    if (username == "" || password == "" || email == "") {
      alert("Please fill in every field");
      return;
    }
    const user = users.find((user) => user[0] === username);
    if (user) {
      alert("Username already taken");
      return;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      alert("Please enter a valid email address");
      return;
    }
    const oldEmail = users.find((oldEmail) => oldEmail[1] === email);
    if (oldEmail) {
      alert("Email already registered");
      return;
    }

    if (password.length < 8) {
      alert("Password must be at least 8 characters");
      return;
    }

    // users.push([username, email, password]);
    console.log("Registration complete!");
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="flex flex-col items-center justify-between h-140 w-100 text-white font-bold p-10 bg-gray-700 rounded-2xl">
        <h1 className="text-3xl pb-6">Register</h1>
        <div className="flex flex-col w-full items-center">
          <div className="w-full m-4">
            <p className="text-sm pb-2">Username</p>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full bg-gray-100 border-4 border-gray-400 text-gray-700 px-4 py-2 text-base rounded-lg outline-none"
              placeholder="Enter username"
            />
          </div>

          <div className="w-full m-4">
            <p className="text-sm pb-2">Email</p>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-gray-100 border-4 border-gray-400 text-gray-700 px-4 py-2 text-base rounded-lg outline-none"
              placeholder="Enter email"
            />
          </div>

          <div className="w-full m-4">
            <p className="text-sm pb-2">Password</p>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-gray-100 border-4 border-gray-400 text-gray-700 px-4 py-2 text-base rounded-lg outline-none"
              placeholder="Enter password"
            />
          </div>

          <button
            id="btn-check"
            onClick={handleRegister}
            className="mt-6 mb-10 w-40 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-700"
          >
            Register
          </button>
        </div>
        <p className="text-color-100 font-normal text-sm">
          Back to
          <Link href="/login" className="ml-2 underline">
            Log in
          </Link>
        </p>
      </div>
    </div>
  );
}
