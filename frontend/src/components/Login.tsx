"use client";
import { database1, database2} from "../Database";

import { useState } from "react";

export default function Login() {
  const [inputValue, setInputValue] = useState("");
  const [passWord, setPassWord] = useState("");
  const [message, setMessage] = useState("");


  const handleLogin = () => {
    // Rendu insensible à la casse pour l'email
    const index = database1.findIndex(
      (email) => email.toLowerCase() === inputValue.toLowerCase()
    );

    if (index !== -1 && database2[index] === passWord) {
      setMessage("✅");
    } else {
      setMessage("❌");
    }
  };

  return (
    <main className="p-6">
      <h1 className="text-2xl font-bold mb-4">Login</h1>

      <input
        id="username"
        type="text"
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        className="border border-gray-400 p-2 rounded"
        placeholder="enter username"
      />

      <input
        id="password"
        type="text"
        value={passWord}
        onChange={(e) => setPassWord(e.target.value)}
        className="border border-gray-400 p-2 rounded ml-2"
        placeholder="enter password"
      />

      <button
        id="btn-check"
        onClick={handleLogin}
        className="ml-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Login
      </button>

      {message && <p className="mt-4">{message}</p>}
    </main>
  );
}
