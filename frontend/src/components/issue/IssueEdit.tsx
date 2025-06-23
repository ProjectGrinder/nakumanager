"use client";

import {
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
} from "@mui/material";
import PriorityIcon from "../PriorityIcon";
import StatusIcon from "../StatusIcon";
import { useState } from "react";

export default function IssueEdit() {
  const issue = {
    name: "Issue 1",
    description: "This is a test issue",
    status: "In Progress",
    priority: "High Priority",
    project: "AI Voicebot",
    creator: "Member 1",
    assignee: "Member 2",
    startDate: "2024-01-01",
    endDate: "2024-06-01",
    label: "Design",
  };
  const [name, setName] = useState(issue.name);
  const [description, setDescription] = useState(issue.description);
  const [status, setStatus] = useState(issue.status);
  const [priority, setPriority] = useState(issue.priority);
  const [assignee, setAssignee] = useState(issue.assignee);
  const [project, setProject] = useState(issue.project);
  const [startDate, setStartDate] = useState(issue.startDate);
  const [endDate, setEndDate] = useState(issue.endDate);
  const [label, setLabel] = useState(issue.label);
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
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex flex-col w-1/2 mb-8">
        <span className="text-2xl font-bold mb-4">Issue Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter issue name"
        />
      </div>
      <div className="flex flex-col gap-4 text-xl font-normal h-60 mb-2">
        <span>Description:</span>
        <FormControl sx={style}>
          <TextField
            id="outlined-multiline-static"
            label="Enter description"
            multiline
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={6}
          />
        </FormControl>
      </div>
      <div className="flex flex-col text-lg font-normal mb-6">
        <div className="flex flex-row justify-between mb-6">
          <div>
            <span>Issue Assignee:</span>
            <FormControl fullWidth sx={style}>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={assignee}
                label="Assignee"
                onChange={(e) => setAssignee(e.target.value)}
              >
                <MenuItem value={"Member 1"}>Member 1</MenuItem>
                <MenuItem value={"Member 2"}>Member 2</MenuItem>
                <MenuItem value={"Member 3"}>Member 3</MenuItem>
              </Select>
            </FormControl>
          </div>
          <div>
            <span>Project:</span>
            <FormControl fullWidth sx={style}>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={project}
                label="Project"
                onChange={(e) => setProject(e.target.value)}
              >
                <MenuItem value={"AI Voicebot"}>AI Voicebot</MenuItem>
                <MenuItem value={"Frontend"}>Frontend</MenuItem>
                <MenuItem value={"Backend"}>Backend</MenuItem>
              </Select>
            </FormControl>
          </div>
        </div>
        <div className="flex flex-row justify-between mb-6">
          <div>
            <span>Status:</span>
            <FormControl fullWidth sx={style}>
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
          <div>
            <span>Priority:</span>
            <FormControl fullWidth sx={style}>
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
        <div className="mt-4">
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
