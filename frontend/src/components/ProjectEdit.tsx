"use client";

import {
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  SelectChangeEvent,
} from "@mui/material";
import PriorityIcon from "./PriorityIcon";
import StatusIcon from "./StatusIcon";
import { useState } from "react";

export default function ProjectEdit() {
  const project = {
    name: "AI Voicebot",
    status: "In Progress",
    priority: "High Priority",
    leader: "Alice",
    startDate: "2024-01-01",
    endDate: "2024-06-01",
    label: "AI",
    members: [
      ["Member 1", "Frontend"],
      ["Member 2", "Frontend"],
      ["Member 3", "Backend"],
    ],
  };
  const [newLeader, setNewLeader] = useState("");
  const [newStatus, setNewStatus] = useState("");
  const [newPriority, setNewPriority] = useState("");
  const addMember = () => {
    console.log("Add new member");
  };
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  const handleLeader = (event: SelectChangeEvent) => {
    setNewLeader(event.target.value as string);
    project.leader = newLeader;
  };
  const handleStatus = (event: SelectChangeEvent) => {
    setNewStatus(event.target.value as string);
    project.status = newStatus;
  };
  const handlePriority = (event: SelectChangeEvent) => {
    setNewPriority(event.target.value as string);
    project.priority = newPriority;
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
      <div className="flex flex-col text-lg font-normal">
        <div>
          <span>Project Leader:</span>
          <FormControl variant="standard">
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={newLeader}
              label="Leader"
              onChange={handleLeader}
            >
              <MenuItem value={"Member 1"}>Member 1</MenuItem>
              <MenuItem value={"Member 2"}>Member 2</MenuItem>
              <MenuItem value={"Member 3"}>Member 3</MenuItem>
            </Select>
          </FormControl>
        </div>
        <div className="flex flex-row justify-between mb-6">
          <div>
            <span>Status:</span>
            <FormControl variant="standard">
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={newStatus}
                label="Status"
                onChange={handleStatus}
              >
                <MenuItem value={"Backlog"}>Backlog</MenuItem>
                <MenuItem value={"Planned"}>Planned</MenuItem>
                <MenuItem value={"In Progress"}>In Progress</MenuItem>
                <MenuItem value={"Planned"}>Completed</MenuItem>
                <MenuItem value={"In Progress"}>Cancelled</MenuItem>
              </Select>
            </FormControl>
          </div>
          <div>
            <span>Priority:</span>
            <FormControl variant="standard">
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={newPriority}
                label="Priority"
                onChange={handlePriority}
              >
                <MenuItem value={"Backlog"}>Backlog</MenuItem>
                <MenuItem value={"Planned"}>Planned</MenuItem>
                <MenuItem value={"In Progress"}>In Progress</MenuItem>
                <MenuItem value={"Planned"}>Completed</MenuItem>
                <MenuItem value={"In Progress"}>Cancelled</MenuItem>
              </Select>
            </FormControl>
          </div>
        </div>
        <div className="flex flex-row justify-between mb-6">
          <div>
            <span>Start Date:</span>
            <input type="date" className="ml-4 color-white" />
          </div>
          <div>
            <span>End Date:</span>
            <input type="date" className="ml-4 color-white" />
          </div>
        </div>
        <div>
          <span className="mb-6">Label:</span>
          <input
            className="bg-gray-100 text-gray-700 p-2 ml-4 text-base rounded outline-none"
            type="text"
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
      <table className="table-auto text-lg w-full text-left text-white mt-6 h-80 max-h-80 overflow-y-auto">
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
