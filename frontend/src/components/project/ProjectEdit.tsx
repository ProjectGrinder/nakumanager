"use client";

import { FormControl, InputLabel, Select, MenuItem } from "@mui/material";
import { useState } from "react";

export default function ProjectEdit() {
  const project = {
    name: "AI Voicebot",
    status: "In Progress",
    priority: "High Priority",
    leader: "Member 1",
    startDate: "2024-01-01",
    endDate: "2024-06-01",
    label: "AI",
    members: [
      ["Member 1", "Frontend"],
      ["Member 2", "Frontend"],
      ["Member 3", "Backend"],
    ],
  };
  const [name, setName] = useState(project.name);
  const [leader, setLeader] = useState(project.leader);
  const [status, setStatus] = useState(project.status);
  const [priority, setPriority] = useState(project.priority);
  const [startDate, setStartDate] = useState(project.startDate);
  const [endDate, setEndDate] = useState(project.endDate);
  const [label, setLabel] = useState(project.label);
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
  const addMember = () => {
    console.log("Add new member");
  };
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex flex-col w-1/2 mb-8">
        <span className="text-2xl font-bold mb-4">Project Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter project name"
        />
      </div>
      <div className="flex flex-col gap-2 text-lg font-normal">
        <div className="flex flew-row w-100 gap-2">
          <span>Project Leader:</span>
          <FormControl fullWidth sx={style}>
            <InputLabel id="demo-simple-select-label">Leader</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={leader}
              label="Leader"
              onChange={(e) => setLeader(e.target.value)}
            >
              <MenuItem value={"Member 1"}>Member 1</MenuItem>
              <MenuItem value={"Member 2"}>Member 2</MenuItem>
              <MenuItem value={"Member 3"}>Member 3</MenuItem>
            </Select>
          </FormControl>
        </div>
        <div className="flex flex-row justify-between">
          <div className="flex flew-row gap-2">
            <span>Status:</span>
            <FormControl fullWidth sx={style}>
              <InputLabel id="demo-simple-select-label">Status</InputLabel>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={status}
                label="Status"
                onChange={(e) => setStatus(e.target.value)}
              >
                <MenuItem value={"Backlog"}>Backlog</MenuItem>
                <MenuItem value={"Planned"}>Planned</MenuItem>
                <MenuItem value={"In Progress"}>In Progress</MenuItem>
                <MenuItem value={"Completed"}>Completed</MenuItem>
                <MenuItem value={"Cancelled"}>Cancelled</MenuItem>
              </Select>
            </FormControl>
          </div>
          <div className="flex flew-row gap-2">
            <span>Priority:</span>
            <FormControl fullWidth sx={style}>
              <InputLabel id="demo-simple-select-label">Priority</InputLabel>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={priority}
                label="Priority"
                onChange={(e) => setPriority(e.target.value)}
              >
                <MenuItem value={"No Priority"}>No Priority</MenuItem>
                <MenuItem value={"Low Priority"}>Low Priority</MenuItem>
                <MenuItem value={"Medium Priority"}>Medium Priority</MenuItem>
                <MenuItem value={"High Priority"}>High Priority</MenuItem>
                <MenuItem value={"Urgent"}>Urgent</MenuItem>
              </Select>
            </FormControl>
          </div>
        </div>
        <div className="flex flex-row justify-between mb-6">
          <div>
            <span>Start Date:</span>
            <input
              className="bg-gray-100 text-gray-700 p-2 ml-4 rounded outline-none"
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
            />
          </div>
          <div>
            <span>End Date:</span>
            <input
              className="bg-gray-100 text-gray-700 p-2 ml-4 rounded outline-none"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
            />
          </div>
        </div>
        <div>
          <span className="mb-6">Label:</span>
          <input
            className="bg-gray-100 text-gray-700 p-2 ml-4 text-base rounded outline-none"
            type="text"
            value={label}
            onChange={(e) => setLabel(e.target.value)}
            placeholder="Enter label"
          />
        </div>
      </div>
      <div className="flex-row text-xl mt-4 font-bold">
        <span>Team Members</span>
        <i
          className="fa-solid fa-square-plus text-2xl ml-8"
          onClick={addMember}
        ></i>
      </div>
      <table className="table-auto text-lg w-full text-left text-white mt-6 h-50 max-h-50 overflow-y-auto">
        <tbody>
          {project.members.map((member, index) => (
            <tr key={index}>
              <td className="w-[50%] max-w-40 p-2 whitespace-nowrap overflow-x-hidden text-ellipsis">
                {member[0]}
              </td>
              <td className="w-[40%] p-2">{member[1]}</td>
            </tr>
          ))}
        </tbody>
      </table>
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
