"use client";
import { database1, database2, database3 } from "../Database";
import { useState } from "react";
export default function Register() {
  //Implement the register functionality here
  const [creatUser, setCreatUser] = useState("");
  const [creatPassWord, setCreatPassWord] = useState("");
  const [email, setEmail] = useState("")
  const [message, setMessage] = useState("");

  const handleRegister = () => {
    //Handle registration logic
    //if contion : setMessage = OK + add to the TAB the new user
    //else : setMessage = NOT OK 
      database1.push(creatUser);
      database2.push(creatPassWord);
      database3.push(email);
      setMessage("✅ Inscription réussie !");
      // Facultatif : réinitialiser les champs
      setCreatUser("");
      setCreatPassWord("");
      setEmail("");
  };

  return (
    <main className="p-6">
      <h1 className="text-2xl font-bold mb-4">Login</h1>

      <input
        id="username"
        type="text"
        value={creatUser}
        onChange={(e) => setCreatUser(e.target.value)}
        className="border border-gray-400 p-2 rounded"
        placeholder="enter username"
      />

      <input
        id="password"
        type="text"
        value={creatPassWord}
        onChange={(e) => setCreatPassWord(e.target.value)}
        className="border border-gray-400 p-2 rounded ml-2"
        placeholder="enter password"
      />

      <input
        id = "email" 
        type="text" 
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        className="border border-gray-400 p-2 rounded ml-2"
        placeholder="enter email"

      />

      <button
        id="btn-check"
        onClick={handleRegister}
        className="ml-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Register
      </button>

      {message && <p className="mt-4">{message}</p>}
    </main>
  );
}
