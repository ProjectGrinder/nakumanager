"use client";

import { Chip } from "@mui/material";

export default function ViewEdit() {
  const view = {
    name: "View 1",
    creator: "Alice Wonder",
    grouping: {
      status: "In Progress",
      priority: "High Priority",
      assignee: "Bob",
      project: "AI Voicebot",
      label: "AI",
      team: "Team 1",
      endDate: "2024-06-01",
    },
    destination: "",
  };
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  const handleGroup = () => {
    //stand in for grouping logic
  };
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex flex-col w-1/2 mb-8">
        <span className="text-2xl font-bold mb-4">Project Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          placeholder="Enter project name"
        />
      </div>
      <div className="flex flex-row">
        <span className="text-xl font-bold mb-4">Group by:</span>
        <Chip label="Status" onClick={handleGroup} />
        <Chip label="Priority" onClick={handleGroup} />
        <Chip label="Assignee" onClick={handleGroup} />
        <Chip label="Project" onClick={handleGroup} />
        <Chip label="Label" onClick={handleGroup} />
        <Chip label="Team" onClick={handleGroup} />
        <Chip label="End Date" onClick={handleGroup} />
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
