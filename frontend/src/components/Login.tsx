"use client";

import Link from "next/link";
import { users } from "../Database";
import { useState } from "react";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    if (username == "" || password == "") {
      alert("Please fill in every field");
      return;
    }
    const user = users.find((user) => user[0] === username);
    if (user) console.log("Login Successful");
    else alert("Username or password is incorrect");
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="flex flex-col items-center justify-between h-140 w-100 text-white font-bold p-10 bg-gray-700 rounded-2xl">
        <h1 className="text-3xl pb-6">Login</h1>
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
            <p className="text-sm pb-2">Password</p>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-gray-100 border-4 border-gray-400 text-gray-700 px-4 py-2 text-base rounded outline-none"
              placeholder="Enter password"
            />
          </div>

          <button
            id="btn-check"
            onClick={handleLogin}
            className="mt-6 mb-10 w-40 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-700"
          >
            Login
          </button>
        </div>
        <p className="text-color-100 font-normal text-sm">
          Not registered yet?
          <Link href="/register" className="ml-2 underline">
            Register here
          </Link>
        </p>
      </div>
    </div>
  );
}
