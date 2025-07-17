"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import CustomAvatar from "./Avatar";

export default function WorkspaceInfo() {
  const owner = ["John Doe", "johndoe@gmail.com"];
  const members = [
    ["Member 1", "member1@gmail.com", "Admin"],
    ["Member 2", "member2@gmail.com", "Admin"],
    ["KLMNOP", "member3@gmail.com", "Member"],
    ["44P", "member4@gmail.com", "Member"],
    ["jane4321", "member5@gmail.com", "Admin"],
  ];
  const [copied, setCopied] = useState(false);
  const router = useRouter();
  const copyInvite = async () => {
    try {
      await navigator.clipboard.writeText("Copy successful!");
      setCopied(true);
      setTimeout(() => setCopied(false), 1000);
    } catch (error) {
      alert("Failed to copy!");
    }
  };
  return (
    <div className="flex flex-col items-start p-6">
      <span className="text-xl text-white font-bold mb-4">
        Workspace Members
      </span>
      <button
        onClick={copyInvite}
        className="px-4 py-2 bg-blue-500 text-sm text-white rounded-md hover:bg-blue-700"
      >
        <i className="fa-solid fa-plus text-xs mr-2"></i>
        Copy invite link
      </button>
      <div className="max-h-160 overflow-y-auto">
        <table className="table-auto border-collapse w-200 text-left mt-4">
          <thead>
            <tr className="h-8 text-xs font-normal text-gray-400 mb-6">
              <td className="w-1/10"></td>
              <td className="w-3/10">Name</td>
              <td className="w-2/5">Email</td>
              <td className="w-1/10">Status</td>
              <td className="w-1/10"></td>
            </tr>
          </thead>
          <tbody>
            <tr className="h-12 text-sm font-normal text-gray-400 hover:bg-gray-800 transition-colors duration-200">
              <td className="rounded-l-md">
                <div className="flex justify-center">
                  <CustomAvatar name={owner[0]} />
                </div>
              </td>
              <td className="text-gray-200 font-medium">{owner[0]}</td>
              <td>{owner[1]}</td>
              <td>Owner</td>
              <td className="rounded-r-md"></td>
            </tr>
            {members.map((member, idx) => (
              <tr
                key={idx}
                className="h-12 text-sm hover:bg-gray-800 text-gray-400 transition-colors duration-200"
              >
                <td className="rounded-l-md">
                  <div className="flex justify-center">
                    <CustomAvatar name={member[0]} />
                  </div>
                </td>
                <td className="text-gray-200 font-medium">{member[0]}</td>
                <td>{member[1]}</td>
                <td>{member[2]}</td>
                <td className="rounded-r-md h-full align-middle">
                  <div className="flex justify-center items-center h-full">
                    <div className="flex justify-center items-center h-6 w-6 rounded-xl cursor-pointer hover:bg-gray-700 transition duration-200">
                      <i className="fa-solid fa-xmark text-gray-500 text-base"></i>
                    </div>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
