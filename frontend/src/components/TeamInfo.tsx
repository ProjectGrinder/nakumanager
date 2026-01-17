"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import CustomAvatar from "./Avatar";
import DeletePopup from "./popup/DeletePopup";
import AddTeamMemberPopup from "./popup/AddTeamMemberPopup";

export default async function TeamInfo() {
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
  const [name, setName] = useState(team.name);
  const canEdit = currentUser === team.creator;
  const router = useRouter();
  const [memberPopup, setMemberPopup] = useState(false);
  const getWokspaceMembers = await fetch(
    "http://localhost:8080/api/workspace/members",
    {
      method: "GET",
    }
  )
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to fetch teams:", err);
    });
  const handleAddMember = async () => {
    const response = await fetch(
      "http://localhost:8080/api/teams/:id/members",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name,
        }),
      }
    );
    if (!response.ok) {
      console.error("Failed to add member");
    } else {
      router.refresh();
    }
  };
  const handleupdate = async () => {
    const response = await fetch("http://localhost:8080/api/teams/:id", {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name,
      }),
    });
    if (!response.ok) {
      console.error("Failed to update project");
    }
  };
  const [del, setDel] = useState(false);
  const handleDelete = async () => {
    const response = await fetch("http://localhost:8080/api/teams/:id", {
      method: "DELETE",
    });
    if (!response.ok) {
      console.error("Failed to delete project");
    }
  };

  return (
    <div className="flex flex-col items-start p-6">
      <div className="flex-row text-white text-xl font-bold mb-4">
        <textarea
          className="resize-none overflow-hidden bg-transparent p-0 leading-snug focus:outline-none"
          rows={1}
          value={name}
          onChange={canEdit ? (e) => setName(e.target.value) : undefined}
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
        {canEdit && (
          <div className="flex flex-row gap-6">
            <button
              className="px-6 py-2 bg-blue-500 text-sm text-white font-base rounded-md hover:bg-blue-700"
              onClick={handleupdate}
            >
              Save
            </button>
            <button
              className="px-6 py-2 bg-red-500 text-sm text-white font-base rounded-md hover:bg-red-700"
              onClick={() => setDel(true)}
            >
              Delete
            </button>
          </div>
        )}
      </div>
      <DeletePopup
        name={team.name}
        open={del}
        onClose={() => setDel(false)}
        onSubmit={handleDelete}
      ></DeletePopup>
      <button
        className="px-4 py-2 bg-blue-500 text-sm text-white rounded-md hover:bg-blue-700"
        disabled={!canEdit}
        onClick={() => setMemberPopup(true)}
      >
        <i className="fa-solid fa-plus text-xs mr-2"></i>
        Add members
      </button>
      <AddTeamMemberPopup
        current={getWokspaceMembers}
        open={memberPopup}
        onClose={() => setMemberPopup(false)}
        onSubmit={handleAddMember}
      />
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
