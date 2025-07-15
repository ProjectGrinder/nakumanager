"use client";

import { useState } from "react";
import DateFormat from "../DateFormat";
import { FormControl, InputLabel, Select, MenuItem } from "@mui/material";
import CustomAvatar from "../Avatar";

export default function ProjectInfo() {
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
  const [name, setName] = useState(project.name);
  const [leader, setLeader] = useState(project.leader);
  const [status, setStatus] = useState(project.status);
  const [priority, setPriority] = useState(project.priority);
  const [startDate, setStartDate] = useState(project.startDate);
  const [endDate, setEndDate] = useState(project.endDate);
  const [label, setLabel] = useState(project.label);
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
  const progress = (total: number, completed: number) => {
    if (completed == 0)
      return <span className="text-base text-red-500">0%</span>;
    const percent = (completed / total) * 100;
    if (percent < 50)
      return (
        <span className="text-base text-red-500">{Math.floor(percent)}%</span>
      );
    if (percent < 75)
      return (
        <span className="text-base text-orange-400">
          {Math.floor(percent)}%
        </span>
      );
    if (percent < 100)
      return (
        <span className="text-base text-yellow-300">
          {Math.floor(percent)}%
        </span>
      );
    if (percent >= 100)
      return <span className="text-base text-green-300">100%</span>;
    return <span className="text-base text-red-500">0%</span>;
  };
  const nameChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    if (e.target.value.length <= 0) {
      return;
    }
    setName(e.target.value);
  };

  return (
    <div className="flex flex-col items-start p-6">
      <div className="flex-row text-white text-xl font-bold mb-4">
        <textarea
          className="w-full resize-none overflow-hidden bg-transparent p-0 text-xl leading-snug focus:outline-none"
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
          <span className="text-base text-gray-400">Leader:</span>
          <FormControl sx={style}>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={leader}
              label="Leader"
              IconComponent={() => null}
              onChange={(e) => setLeader(e.target.value)}
              MenuProps={menuStyle}
            >
              <MenuItem value={"Not Assigned"}>Not Assigned</MenuItem>
              {project.members.map((member) => (
                <MenuItem value={member[0]}>{member[0]}</MenuItem>
              ))}
            </Select>
          </FormControl>
        </div>
        <hr className="border-gray-500 mt-4 mb-4"></hr>
        <div className="flex flex-row gap-4 items-center mb-4">
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
            <input
              className="bg-gray-700 text-gray-200 text-sm px-2 py-1.5 rounded-lg outline-none appearance-none"
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
            />
            <span className="mx-2 text-base text-gray-200">to</span>
            <input
              className="bg-gray-700 text-gray-200 text-sm px-2 py-1.5 rounded-lg outline-none appearance-none"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
            />
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
        <hr className="border-gray-500 mt-4 mb-4"></hr>
        <div className="flex flex-row gap-16 text-base text-gray-400">
          <div className="flex flex-row gap-2">
            <span>Total issues:</span>
            <span>0</span>
          </div>
          <div className="flex flex-row gap-2">
            <span>Completed issues:</span>
            <span>0</span>
          </div>
          <div className=" flex flex-row gap-2">
            <span>Progress:</span>
            <span>{progress(0, 0)}</span>
          </div>
        </div>
        <hr className="border-gray-500 mt-4 mb-4"></hr>
      </div>
      <div className="flex flex-col">
        <span className="text-lg font-semibold text-gray-200 mb-4">
          Team Members
        </span>
        <button className="px-4 py-2 bg-blue-500 text-sm text-white rounded-md hover:bg-blue-700">
          <i className="fa-solid fa-plus text-xs mr-2"></i>
          Add members
        </button>
      </div>
      <div className="max-h-90 overflow-y-auto">
        <table className="w-200 text-left mt-4">
          <thead>
            <tr className="h-8 text-xs font-normal text-gray-400 mb-6">
              <td className="w-1/10"></td>
              <td className="w-2/5">Name</td>
              <td className="w-2/5">Role</td>
              <td className="w-1/10"></td>
            </tr>
          </thead>
          <tbody>
            {project.members.map((member, idx) => (
              <tr
                key={idx}
                className="h-12 text-sm hover:bg-gray-800 text-gray-400 transition-colors"
              >
                <td className="rounded-l-md">
                  <div className="flex justify-center">
                    <CustomAvatar name={member[0]} />
                  </div>
                </td>
                <td className="text-gray-200 font-medium">{member[0]}</td>
                <td>{member[1]}</td>
                <td className="rounded-r-md h-full align-middle">
                  <div className="flex justify-center items-center h-full">
                    <div className="flex justify-center items-center h-6 w-6 rounded-xl cursor-pointer hover:bg-gray-700 transition duration-200">
                      <i className="fa-solid fa-xmark text-gray-500 text-base"></i>
                    </div>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
