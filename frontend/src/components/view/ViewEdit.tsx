"use client";

import {
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
} from "@mui/material";
import { useState } from "react";
import IssueSelectItem from "../issue/IssueSelectItem";

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
  const issue_list = [
    [
      "Issue 1",
      "In Progress",
      "Urgent",
      "Alice",
      "2024-01-01",
      "2024-06-01",
      "Issue 1",
    ],
    [
      "Frontend",
      "Completed",
      "Low Priority",
      "Bob",
      "2024-02-01",
      "2024-07-01",
      "Issue 2",
    ],
  ];
  const [name, setName] = useState(view.name);
  const [status, setStatus] = useState(view.grouping.status);
  const [priority, setPriority] = useState(view.grouping.priority);
  const [assignee, setAssignee] = useState(view.grouping.assignee);
  const [project, setProject] = useState(view.grouping.project);
  const [label, setLabel] = useState(view.grouping.label);
  const [team, setTeam] = useState(view.grouping.team);
  const [endDate, setEndDate] = useState(view.grouping.endDate);
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
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
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex flex-col w-1/2 mb-8">
        <span className="text-2xl font-bold mb-4">View Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter project name"
        />
      </div>
      <div className="flex flex-col w-200 mb-6">
        <span className="text-xl font-bold mb-8">Group by:</span>
        <div className="flex flex-row gap-4">
          <FormControl fullWidth sx={style}>
            <InputLabel id="demo-simple-select-label">Status</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={status}
              label="Status"
              onChange={(e) => setStatus(e.target.value)}
            >
              <MenuItem value={""}>Not Selected</MenuItem>
              <MenuItem value={"Backlog"}>Backlog</MenuItem>
              <MenuItem value={"Planned"}>Planned</MenuItem>
              <MenuItem value={"In Progress"}>In Progress</MenuItem>
              <MenuItem value={"Completed"}>Completed</MenuItem>
              <MenuItem value={"Cancelled"}>Cancelled</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth sx={style}>
            <InputLabel id="demo-simple-select-label">Priority</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={priority}
              label="Priority"
              onChange={(e) => setPriority(e.target.value)}
            >
              <MenuItem value={""}>Not Selected</MenuItem>
              <MenuItem value={"No Priority"}>No Priority</MenuItem>
              <MenuItem value={"Low Priority"}>Low Priority</MenuItem>
              <MenuItem value={"Medium Priority"}>Medium Priority</MenuItem>
              <MenuItem value={"High Priority"}>High Priority</MenuItem>
              <MenuItem value={"Urgent"}>Urgent</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth sx={style}>
            <InputLabel id="demo-simple-select-label">Assignee</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={assignee}
              label="Assignee"
              onChange={(e) => setAssignee(e.target.value)}
            >
              <MenuItem value={""}>Not Selected</MenuItem>
              <MenuItem value={"Alice"}>Alice</MenuItem>
              <MenuItem value={"Bob"}>Bob</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth sx={style}>
            <InputLabel id="demo-simple-select-label">Project</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={project}
              label="Project"
              onChange={(e) => setProject(e.target.value)}
            >
              <MenuItem value={""}>Not Selected</MenuItem>
              <MenuItem value={"AI Voicebot"}>AI Voicebot</MenuItem>
              <MenuItem value={"Frontend"}>Frontend</MenuItem>
              <MenuItem value={"Backend"}>Backend</MenuItem>
            </Select>
          </FormControl>
        </div>
        <div className="flex flex-row mt-8 gap-4 items-center">
          <FormControl fullWidth sx={style}>
            <TextField
              id="outlined-basic"
              label="Label"
              variant="outlined"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
            />
          </FormControl>
          <FormControl fullWidth sx={style}>
            <InputLabel id="demo-simple-select-label">Team</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={team}
              label="Team"
              onChange={(e) => setTeam(e.target.value)}
            >
              <MenuItem value={""}>Not Selected</MenuItem>
              <MenuItem value={"Team 1"}>Team 1</MenuItem>
              <MenuItem value={"Team 2"}>Team 2</MenuItem>
              <MenuItem value={"Team 3"}>Team 3</MenuItem>
            </Select>
          </FormControl>
          <div>
            <label className="text-sm">End Date</label>
            <input
              className="bg-gray-100 text-gray-700 p-2 rounded outline-none"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
            />
          </div>
        </div>
      </div>
      <div className="w-300 overflow-y-auto">
        {issue_list.map((issue, index) => (
          <IssueSelectItem
            name={issue[0]}
            status={issue[1]}
            priority={issue[2]}
            assigned={issue[3]}
            startDate={issue[4]}
            endDate={issue[5]}
            destination={issue[6]}
            key={index}
          />
        ))}
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
