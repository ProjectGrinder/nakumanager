"use client";

import { useState } from "react";

export default function TeamEdit() {
  const team_members = [
    ["Admin 1", "Project Manager"],
    ["Admin 2", "Designer"],
    ["Member 1", "Frontend"],
    ["Member 2", "Backend"],
    [
      "Member 3",
      "A very long role name that should be truncated if it exceeds the width of the table cell",
    ],
  ]; //Role changing has issue
  const [value, setValue] = useState("");
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  const addMember = () => {
    console.log("Add member clicked");
  };
  const removeMember = (index: number) => {
    team_members.splice(index, 1);
    console.log(`Member at index ${index} removed`);
  };
  return (
    <div className="relative flex flex-col p-10 text-white w-7/10">
      <div className="flex flex-col w-1/2 mb-16">
        <span className="text-2xl font-bold mb-4">Team Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          placeholder="Enter team name"
        />
      </div>
      <div className="flex-column items-left justify-start w-1/2">
        <div className="flex flex-row items-center">
          <span className="font-bold text-2xl">Team Members</span>
          <i
            className="fa-solid fa-square-plus text-2xl ml-8"
            onClick={addMember}
          ></i>
        </div>
        <table className="table-auto text-lg w-full text-left text-white mt-6 h-80 max-h-80 overflow-y-auto">
          <tbody>
            {team_members.map((member, index) => (
              <tr>
                <td
                  className="w-[50%] max-w-40 p-4 whitespace-nowrap overflow-x-hidden text-ellipsis"
                  key={index}
                >
                  {member[0]}
                </td>
                <td className="w-[40%] p-4" key={index}>
                  <input
                    className="h-8 w-full text-base bg-gray-100 border-gray-500 text-gray-700 rounded outline-none p-2"
                    type="text"
                    value={member[1]}
                    onChange={(e) => setValue(e.target.value)}
                  />
                </td>
                <td>
                  <i
                    className="fa-solid fa-ban ml-4"
                    onClick={() => removeMember(index)}
                  ></i>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="flex flex-row fixed bottom-20 right-40">
        <button className="text-gray-100 px-8 py-2" onClick={handleCancel}>
          Cancel
        </button>
        <button
          className="bg-gray-100 text-gray-700 px-8 py-2 rounded hover:bg-gray-300"
          onClick={handleSave}
        >
          Save
        </button>
      </div>
    </div>
  );
}
