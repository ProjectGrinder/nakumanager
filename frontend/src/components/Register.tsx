"use client";
import { useState } from "react";
export default function Register() {
  //Implement the register functionality here
  const [creatUser, setCreatUser] = useState("");
  const [creatPassWord, setCreatPassWord] = useState("");
  const [email, setEmail] = useState("")
  const [message, setMessage] = useState("");

  const handleRegister = () => {
  // Récupère les anciens utilisateurs du localStorage (ou [] s'il n'y en a pas)
  const users = JSON.parse(localStorage.getItem("users") || "[]");

  // Crée un nouvel utilisateur
  const newUser = {
    username: creatUser,
    password: creatPassWord,
    email: email,
  };

  // Ajoute et sauvegarde dans localStorage
  users.push(newUser);
  localStorage.setItem("users", JSON.stringify(users));
  console.log("Utilisateurs enregistrés :", users);

  setMessage("✅ successful register!");
  setCreatUser("");
  setCreatPassWord("");
  setEmail("");
};


  return (
  <main className="min-h-screen flex items-center justify-center bg-[#2f2f2f]">
  <div className="bg-[#5a5a5a] p-10 rounded-xl shadow-lg w-[350px] text-white">
    <h1 className="text-3xl font-bold text-center mb-6">Register</h1>
    
    <label htmlFor="username" className="block font-semibold mb-1">Username</label>
    <input
      id="username"
      type="text"
      value={creatUser} // ou creatUser
      onChange={(e) => setCreatUser(e.target.value)}
      className="w-full p-2 rounded mb-4 text-black bg-white"
      placeholder="Enter username"
    />

    {/* Email field only on Register page */}
    <label htmlFor="email" className="block font-semibold mb-1">Email</label>
    <input
      id="email"
      type="text"
      value={email}
      onChange={(e) => setEmail(e.target.value)}
      className="w-full p-2 rounded mb-4 text-black bg-white"
      placeholder="Enter email"
    />

    <label htmlFor="password" className="block font-semibold mb-1">Password</label>
    <input
      id="password"
      type="password"
      value={creatPassWord} // ou creatPassWord
      onChange={(e) => setCreatPassWord(e.target.value)}
      className="w-full p-2 rounded mb-6 text-black bg-white"
      placeholder="Enter password"
    />

    <button
      onClick={handleRegister} // ou handleRegister
      className="w-full bg-white text-black font-semibold py-2 rounded hover:bg-gray-200 transition"
    >
      Register 
    </button>


 
    <div className="text-center mt-4 text-sm">
      Back to{" "}
      <a href="/login" className="underline">
        Log In
      </a>
    </div>

    {message && <p className="mt-4 text-center">{message}</p>}
  </div>
</main>

  );
}
