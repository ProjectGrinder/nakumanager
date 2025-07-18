"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import CustomAvatar from "./Avatar";

export default function TeamInfo() {
  const currentUser = "John Doe";
  const team = {
    name: "Team 1",
    creator: "John Doe",
    members: [
      ["John Doe", "Manager"],
      ["Member 1", "Data Scientist"],
      ["Member 2", "Designer"],
      ["KLMNOP", "Frontend"],
      ["44P", "Backend"],
      ["jane4321", "Tester"],
    ],
  };
  const [teamName, setTeamName] = useState(team.name);
  const canEdit = currentUser === team.creator;
  const router = useRouter();
  return (
    <div className="flex flex-col items-start p-6">
      <div className="flex-row text-white text-xl font-bold mb-4">
        <textarea
          className="resize-none overflow-hidden bg-transparent p-0 leading-snug focus:outline-none"
          rows={1}
          value={teamName}
          onChange={canEdit ? (e) => setTeamName(e.target.value) : undefined}
          onInput={(e) => {
            const textarea = e.currentTarget;
            textarea.style.height = "auto";
            textarea.style.width = "auto";
            textarea.style.height = textarea.scrollHeight + "px";
            textarea.style.width = textarea.scrollWidth + "px";
          }}
          spellCheck={false}
          autoCorrect="off"
          autoCapitalize="off"
        ></textarea>
      </div>
      <button
        className="px-4 py-2 bg-blue-500 text-sm text-white rounded-md hover:bg-blue-700"
        disabled={!canEdit}
      >
        <i className="fa-solid fa-plus text-xs mr-2"></i>
        Add members
      </button>
      <div className="max-h-160 overflow-y-auto">
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
            {team.members.map((member, idx) => (
              <tr
                key={idx}
                className="h-12 text-sm hover:bg-gray-800 text-gray-400 transition-colors"
              >
                <td className="rounded-l-md">
                  <div className="flex justify-center">
                    <CustomAvatar name={member[0]} />
                  </div>
                </td>
                <td className="text-gray-200 font-medium">{member[0]}</td>
                <td>{member[1]}</td>
                <td className="rounded-r-md h-full align-middle">
                  <div className="flex justify-center items-center h-full">
                    {canEdit && (
                      <div className="flex justify-center items-center h-6 w-6 rounded-xl cursor-pointer hover:bg-gray-700 transition duration-200">
                        <i className="fa-solid fa-xmark text-gray-500 text-base"></i>
                      </div>
                    )}
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
