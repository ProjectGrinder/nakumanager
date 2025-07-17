"use client";

import { useState } from "react";
import { SelectChangeEvent } from "@mui/material/Select";
import { FormControlLabel, FormGroup, Switch } from "@mui/material";

export default function WorkspaceEdit() {
  const workspaceName = "Workspace 1";
  const member_info = [
    ["Admin 1", "Admin"],
    ["Admin 2", "Admin"],
    ["Admin 3", "Admin"],
    ["Member 1", "Member"],
    ["Member 2", "Member"],
    ["Member 3", "Member"],
  ];
  const [info, setInfo] = useState(member_info);
  const [name, setName] = useState(workspaceName);
  const [copied, setCopied] = useState(false);
  const sendInvite = async () => {
    try {
      await navigator.clipboard.writeText("Copy successful!");
      setCopied(true);
      setTimeout(() => setCopied(false), 1000);
    } catch (error) {
      alert("Failed to copy!");
    }
  };
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  const removeMember = (index: number) => {
    member_info.splice(index, 1);
    console.log(`Member at index ${index} removed`);
  };
  return (
    <div className="relative flex flex-col p-10 text-white w-7/10">
      <div className="flex flex-col w-1/2 mb-16">
        <span className="text-2xl font-bold mb-4">Workspace Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter workspace name"
        />
      </div>
      <div className="flex-column items-left justify-start w-1/2">
        <div className="flex flex-row items-center">
          <span className="font-bold text-2xl">Members</span>
          <button
            className="ml-12 bg-gray-200 text-gray-700 text-sm px-3 py-2 rounded"
            onClick={sendInvite}
          >
            + Invite Link
          </button>
        </div>
        <table className="table-auto text-lg w-full text-left text-white mt-6 h-80 overflow-y-auto">
          <tbody>
            {info.map((member, index) => (
              <tr key={index}>
                <td className="w-[60%] p-3">{member[0]}</td>
                <td className="w-[30%] p-3">
                  <FormGroup>
                    <FormControlLabel
                      control={<Switch checked={member[1] === "Admin"} />}
                      onChange={() => {
                        const updated = info.map((m, i) =>
                          i === index
                            ? [m[0], m[1] === "Admin" ? "Member" : "Admin"]
                            : m
                        );
                        setInfo(updated);
                      }}
                      label="Admin"
                    />
                  </FormGroup>
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
