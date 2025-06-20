"use client";

import {
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  SelectChangeEvent,
  TextField,
} from "@mui/material";
import PriorityIcon from "./PriorityIcon";
import StatusIcon from "./StatusIcon";
import { useState } from "react";

export default function IssueEdit() {
  const issue = {
    name: "Issue 1",
    description: "This is a test issue",
    status: "In Progress",
    priority: "High Priority",
    project: "AI Voicebot",
    creator: "Alice",
    assignee: "Bob",
    startDate: "2024-01-01",
    endDate: "2024-06-01",
    label: "Design",
  };
  const [newAssignee, setNewAssignee] = useState("");
  const [newStatus, setNewStatus] = useState("");
  const [newPriority, setNewPriority] = useState("");
  const handleCancel = () => {
    console.log("Cancel clicked");
  };
  const handleSave = () => {
    console.log("Save clicked");
  };
  const handleAssignee = (event: SelectChangeEvent) => {
    setNewAssignee(event.target.value as string);
    issue.assignee = newAssignee;
  };
  const handleStatus = (event: SelectChangeEvent) => {
    setNewStatus(event.target.value as string);
    issue.status = newStatus;
  };
  const handlePriority = (event: SelectChangeEvent) => {
    setNewPriority(event.target.value as string);
    issue.priority = newPriority;
  };
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex flex-col w-1/2 mb-8">
        <span className="text-2xl font-bold mb-4">Issue Name</span>
        <input
          className="bg-gray-100 text-gray-700 p-4 text-base rounded outline-none"
          type="text"
          placeholder="Enter issue name"
        />
      </div>
      <div className="flex flex-col gap-4 text-xl font-normal h-60 mb-8">
        <span>Description:</span>
        <TextField
          id="outlined-multiline-static"
          label="Enter description"
          multiline
          rows={6}
        />
      </div>
      <div className="flex flex-col text-lg font-normal mb-6">
        <div>
          <span>Issue Assignee:</span>
          <FormControl variant="standard">
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={newAssignee}
              label="Leader"
              onChange={handleAssignee}
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
