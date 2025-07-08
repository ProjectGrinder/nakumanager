"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import CustomAvatar from "./Avatar";

export default function WorkspaceInfo() {
  const owner = "John Doe";
  const members = [
    ["John Doe", "Manager"],
    ["Member 1", "Data Scientist"],
    ["Member 2", "Designer"],
    ["KLMNOP", "Frontend"],
    ["44P", "Backend"],
    ["jane4321", "Tester"],
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
    <div className="flex flex-col items-start p-6 text-white">
      <span className="text-xl font-bold mb-4">Team Members</span>
      {/* <i
          className="fa-solid fa-gear text-2xl ml-10"
          onClick={editWorkspace}
        ></i> */}
      <button className="px-4 py-2 bg-blue-500 text-sm text-white rounded-md hover:bg-blue-700">
        <i className="fa-solid fa-plus text-xs mr-2"></i>
        Add members
      </button>
      <table className="w-200 text-left mt-4">
        <thead>
          <tr className="h-8 text-xs font-normal text-gray-400 mb-6">
            <td className="w-1/10"></td>
            <td className="w-2/5">Name</td>
            <td className="w-2/5">Role</td>
            <td className="w-1/10"></td>
          </tr>
        </thead>
        <tbody>
          {members.map((member, idx) => (
            <tr
              key={idx}
              className="h-12 text-sm hover:bg-gray-800 text-gray-400 transition-colors"
            >
              <td>
                <div className="flex justify-center">
                  <CustomAvatar name={member[0]} />
                </div>
              </td>
              <td className="text-gray-200 font-medium">{member[0]}</td>
              <td>{member[1]}</td>
              <td></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
