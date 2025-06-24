"use client";
import { useState } from "react";

export default function Login() {
  const [inputValue, setInputValue] = useState("");
  const [passWord, setPassWord] = useState("");
  const [message, setMessage] = useState("");

  const handleLogin = () => {
    const users = JSON.parse(localStorage.getItem("users") || "[]");

  const user = users.find(
    (u: { username: string; password: string }) =>
      u.username.toLowerCase() === inputValue.toLowerCase() &&
      u.password === passWord
  );

    if (user) {
      setMessage("✅ Connexion succesful !");
    } else {
      console.log(users);
      setMessage("❌ incorrects");
    }
  };

  return (
<main className="min-h-screen flex items-center justify-center bg-[#2f2f2f]">
  <div className="bg-[#5a5a5a] p-10 rounded-xl shadow-lg w-[350px] text-white">
    <h1 className="text-3xl font-bold text-center mb-6">Log In</h1> {/* ou "Register" */}
    
    <label htmlFor="username" className="block font-semibold mb-1">Username</label>
    <input
      id="username"
      type="text"
      value={inputValue} // ou creatUser
      onChange={(e) => setInputValue(e.target.value)}
      className="w-full p-2 rounded mb-4 text-black bg-white"
      placeholder="Enter username"
    />


    <label htmlFor="password" className="block font-semibold mb-1">Password</label>
    <input
      id="password"
      type="password"
      value={passWord} // ou creatPassWord
      onChange={(e) => setPassWord(e.target.value)}
      className="w-full p-2 rounded mb-6 text-black bg-white"
      placeholder="Enter password"
    />

    <button
      onClick={handleLogin} // ou handleRegister
      className="w-full bg-white text-black font-semibold py-2 rounded hover:bg-gray-200 transition"
    >
      Log In {/* ou Register */}
    </button>

    <div className="text-center mt-4 text-sm">
      Not registered yet?{" "}
      <a href="/register" className="underline">
        Click here
      </a>
    </div>

    {message && <p className="mt-4 text-center">{message}</p>}
  </div>
</main>

  );
}
