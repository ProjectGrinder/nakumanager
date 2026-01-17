"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
export default function Register() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [email, setEmail] = useState("");

  const [message, setMessage] = useState("");
  const router = useRouter();

  const handleRegister = async () => {
    if (username == "" || password == "" || email == "") {
      alert("Please fill in every field");
      return;
    }
    try {
      const res = await fetch("http://localhost:8080/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`API error: ${res.status} - ${errorText}`);
      }

      const data = await res.json();
      setMessage(data.message);
      router.push("/");
    } catch (err) {
      console.error(err);
      setMessage("Registration failed");
    }
    console.log(message);
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
        <div className="flex flex-col w-full items-center mb-15">
          <span className="text-2xl font-semibold text-gray-200">Register</span>
          <div className="w-full m-2">
            <p className="text-sm pb-2 text-gray-200">Username</p>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full bg-gray-200 border-2 border-gray-400 text-gray-700 px-4 py-2 text-sm rounded-lg outline-none"
            />
          </div>

          <div className="w-full m-2">
            <p className="text-sm pb-2 text-gray-200">Email</p>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-gray-200 border-2 border-gray-400 text-gray-700 px-4 py-2 text-sm rounded-lg outline-none"
            />
          </div>

          <div className="w-full m-2">
            <p className="text-sm pb-2 text-gray-200">Password</p>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-gray-200 border-2 border-gray-400 text-gray-700 px-4 py-2 text-sm rounded-lg outline-none"
            />
          </div>

          <button
            id="btn-check"
            onClick={handleRegister}
            className="mt-6 border-box px-10 py-4 bg-blue-500 text-lg text-white rounded-xl hover:bg-blue-700"
          >
            Confirm
          </button>
        </div>
        <p className="text-gray-100 font-normal text-sm">
          Back to
          <Link href="/" className="ml-2 underline">
            Log in
          </Link>
        </p>
      </div>
    </div>
  );
}
