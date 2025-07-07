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
      <div className="flex flex-col items-center justify-between h-160 w-120 text-white font-bold p-10 bg-gray-700 rounded-2xl">
        <div className="flex flex-col items-center">
          <h1 className="text-3xl text-gray-200 pb-4">Nakumanager</h1>
          <p className="text-sm text-gray-400 pb-6">
            The management software that grows with your need
          </p>
        </div>
        <div className="flex flex-col w-full items-center mb-30">
          <span className="text-2xl font-semibold text-gray-200">Login</span>
          <div className="w-full m-2">
            <p className="text-sm pb-2 text-gray-200">Username</p>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full bg-gray-200 border-2 border-gray-400 text-gray-700 px-4 py-2 text-sm rounded-xl outline-none"
            />
          </div>

          <div className="w-full m-2">
            <p className="text-sm pb-2 text-gray-200">Password</p>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-gray-200 border-2 border-gray-400 text-gray-700 px-4 py-2 text-sm rounded-xl outline-none"
            />
          </div>

          <button
            id="btn-check"
            onClick={handleLogin}
            className="mt-6 border-box px-10 py-4 bg-blue-500 text-lg text-white rounded-xl hover:bg-blue-700"
          >
            Confirm
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
