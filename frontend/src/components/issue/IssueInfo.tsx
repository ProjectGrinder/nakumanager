"use client";

import { useState } from "react";
import { FormControl, Select, MenuItem } from "@mui/material";
import CustomDatePicker from "../CustomDatePicker";

export default function IssueInfo() {
  const issue = {
    name: "Issue 1",
    description: "This is a test issue",
    status: "Planned",
    priority: "High Priority",
    project: "AI Voicebot",
    creator: "Member 1",
    assignee: "Member 2",
    startDate: new Date("2024-01-01"),
    endDate: new Date("2024-06-01"),
    label: "Design",
  };
  const projects = ["AI Voicebot", "Web App", "Mobile App"];
  const members = ["Member 1", "Member 2", "Member 3"];
  const [name, setName] = useState(issue.name);
  const [description, setDescription] = useState(issue.description);
  const [status, setStatus] = useState(issue.status);
  const [priority, setPriority] = useState(issue.priority);
  const [assignee, setAssignee] = useState(issue.assignee);
  const [project, setProject] = useState(issue.project);
  const [startDate, setStartDate] = useState<Date | null>(issue.startDate);
  const [endDate, setEndDate] = useState<Date | null>(issue.endDate);
  const [label, setLabel] = useState(issue.label);
  const style = {
    width: "auto",
    border: "none",
    boxShadow: "none",
    "& .MuiInputLabel-root": { color: "#e5e7eb" },
    "& .MuiSelect-select": {
      padding: "4px 10px",
      minHeight: "unset",
      fontSize: "0.9rem",
      display: "flex",
      alignItems: "center",
      paddingRight: "10px !important",
      minWidth: "fit-content",
    },
    "& .MuiOutlinedInput-root": {
      color: "#e5e7eb",
      backgroundColor: "#374151",
      borderRadius: "0.5rem",
      "& fieldset": {
        border: "none",
      },
      "&:hover fieldset": {
        borderColor: "#6b7280",
      },
      "&.Mui-focused fieldset": {
        borderColor: "#9ca3af",
      },
    },
  };
  const menuStyle = {
    PaperProps: {
      sx: {
        backgroundColor: "#374151",
        color: "#e5e7eb",
        borderRadius: 2,
        "& .MuiMenuItem-root": {
          fontSize: "0.9rem",
        },
      },
    },
  };

  const nameChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    if (e.target.value.length <= 0) {
      return;
    }
    setName(e.target.value);
  };

  return (
    <div className="flex flex-col items-start p-6 w-150">
      <div className="flex-row text-white text-xl font-bold mb-4">
        <textarea
          className="min-w-[30rem] resize-none overflow-hidden bg-transparent p-0 leading-snug focus:outline-none"
          rows={1}
          value={name}
          onChange={nameChange}
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
      </div>
      <div className="flex flex-col font-normal">
        <div className="flex flex-row gap-2 items-center">
          <span className="text-base text-gray-400">
            Creator: {issue.creator}
          </span>
        </div>
        <textarea
          className="w-150 bg-gray-700 text-gray-200 text-sm p-4 rounded-lg my-4 focus:outline-none"
          rows={5}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="No description..."
        ></textarea>
        <hr className="border-gray-500 mt-4 mb-4"></hr>
        <div className="flex flex-row gap-6 items-center mb-4">
          <FormControl fullWidth sx={style}>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={status}
              label="Status"
              IconComponent={() => null}
              onChange={(e) => setStatus(e.target.value)}
              MenuProps={menuStyle}
            >
              <MenuItem value={"Backlog"}>
                <i className="fa-solid fa-circle-notch text-gray-200 text-sm mr-2"></i>
                Backlog
              </MenuItem>
              <MenuItem value={"Planned"}>
                <i className="fa-solid fa-pen-ruler text-sky-300 text-sm mr-2"></i>
                Planned
              </MenuItem>
              <MenuItem value={"In Progress"}>
                <i className="fa-solid fa-clock text-yellow-300 text-sm mr-2"></i>
                In Progress
              </MenuItem>
              <MenuItem value={"Completed"}>
                <i className="fa-solid fa-circle-check text-green-300 text-sm mr-2"></i>
                Completed
              </MenuItem>
              <MenuItem value={"Cancelled"}>
                <i className="fa-solid fa-circle-xmark text-red-500 text-sm mr-2"></i>
                Cancelled
              </MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth sx={style}>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={priority}
              label="Priority"
              IconComponent={() => null}
              onChange={(e) => setPriority(e.target.value)}
              MenuProps={menuStyle}
            >
              <MenuItem value={"No Priority"}>
                <i className="fa-solid fa-circle-exclamation text-gray-200 text-sm mr-2"></i>
                No Priority
              </MenuItem>
              <MenuItem value={"Low Priority"}>
                <i className="fa-solid fa-circle-exclamation text-green-300 text-sm mr-2"></i>
                Low Priority
              </MenuItem>
              <MenuItem value={"Medium Priority"}>
                <i className="fa-solid fa-circle-exclamation text-yellow-300 text-sm mr-2"></i>
                Medium Priority
              </MenuItem>
              <MenuItem value={"High Priority"}>
                <i className="fa-solid fa-circle-exclamation text-orange-400 text-sm mr-2"></i>
                High Priority
              </MenuItem>
              <MenuItem value={"Urgent"}>
                <i className="fa-solid fa-circle-exclamation text-red-500 text-sm mr-2"></i>
                Urgent
              </MenuItem>
            </Select>
          </FormControl>
          <div>
            <CustomDatePicker value={startDate} onChange={setStartDate} />
            <span className="mx-2 text-base text-gray-200">to</span>
            <CustomDatePicker value={endDate} onChange={setEndDate} />
          </div>
        </div>
        <div className="flex flex-row gap-8 mb-4">
          <div>
            <FormControl sx={style}>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={assignee}
                label="Assignee"
                IconComponent={() => null}
                onChange={(e) => setAssignee(e.target.value)}
                MenuProps={menuStyle}
              >
                <MenuItem value={"Not Assigned"}>
                  <i className="fa-solid fa-circle-user text-gray-400 mr-2"></i>
                  Not Assigned
                </MenuItem>
                {members.map((member) => (
                  <MenuItem value={member}>
                    <i className="fa-solid fa-circle-user text-gray-200 mr-2"></i>
                    {member}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </div>
          <div>
            <FormControl sx={style}>
              <Select
                labelId="demo-simple-select-label"
                id="demo-simple-select"
                value={project}
                label="Project"
                IconComponent={() => null}
                onChange={(e) => setProject(e.target.value)}
                MenuProps={menuStyle}
              >
                <MenuItem value={"Not Assigned"}>
                  <i className="fa-solid fa-cube text-gray-400 mr-2"></i>Not
                  Assigned
                </MenuItem>
                {projects.map((project) => (
                  <MenuItem value={project}>
                    <i className="fa-solid fa-cube text-gray-200 mr-2"></i>
                    {project}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </div>
        </div>
        <div className="flex flex-row gap-2 items-center">
          <span className="text-base text-gray-400">Label:</span>
          <input
            className="bg-gray-700 text-gray-200 p-2 text-sm rounded-lg w-40 h-8 outline-none"
            type="text"
            value={label}
            onChange={(e) => setLabel(e.target.value)}
            placeholder="None"
          />
        </div>
      </div>
    </div>
  );
}
