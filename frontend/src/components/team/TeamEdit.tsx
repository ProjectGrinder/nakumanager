"use client";

import { FormControl, TextField } from "@mui/material";
import { useState } from "react";

export default function TeamEdit() {
  const team = [
    ["Admin 1", "Project Manager"],
    ["Admin 2", "Designer"],
    ["Member 1", "Frontend"],
    ["Member 2", "Backend"],
    [
      "Member 3",
      "A very long role name that should be truncated if it exceeds the width of the table cell",
    ],
  ];
  const [teamList, setTeamList] = useState(team);
  const style = {
    "& .MuiInputLabel-root": { color: "white" },
    "& .MuiOutlinedInput-root": {
      color: "white",
      backgroundColor: "#374151",
      borderRadius: "0.5rem",
      "& fieldset": {
        borderColor: "#d1d5db",
      },
      "&:hover fieldset": {
        borderColor: "#60a5fa",
      },
      "&.Mui-focused fieldset": {
        borderColor: "#2563eb",
      },
    },
    "& .MuiSvgIcon-root": {
      color: "white",
    },
  };
  const handleRoleChange = (index: number, value: string) => {
    setTeamList((prev) =>
      prev.map((member, i) => (i === index ? [member[0], value] : member))
    );
  };
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
    setTeamList((prev) => prev.filter((_, i) => i !== index));
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
            {teamList.map((member, index) => (
              <tr key={index}>
                <td className="w-[50%] max-w-40 p-3 whitespace-nowrap overflow-x-hidden text-ellipsis">
                  {member[0]}
                </td>
                <td className="w-[40%] p-3">
                  <FormControl fullWidth sx={style}>
                    <TextField
                      id="outlined-basic"
                      label="Label"
                      variant="outlined"
                      value={member[1]}
                      onChange={(e) => handleRoleChange(index, e.target.value)}
                    />
                  </FormControl>
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
