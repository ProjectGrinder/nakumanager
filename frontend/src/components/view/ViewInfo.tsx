"use client";

import { FormControl, Select, MenuItem } from "@mui/material";
import CustomDatePicker from "../CustomDatePicker";
import { useState } from "react";
import IssueSelectItem from "../issue/IssueSelectItem";

export default function ViewInfo() {
  const view = {
    name: "View 1",
    creator: "Member 1",
    status: "Not set",
    priority: "Not set",
    assignee: "Not set",
    team: "Not set",
    project: "Not set",
    label: "Not set",
    endDate: new Date("2024-06-01"),
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
  const projects = ["AI Voicebot", "Web App", "Mobile App"];
  const members = ["Member 1", "Member 2", "Member 3"];
  const teams = ["Team 1", "Team 2", "Team 69"];
  const labels = ["AI", "Design", "BOTNOI"];
  const [name, setName] = useState(view.name);
  const [status, setStatus] = useState(view.status);
  const [priority, setPriority] = useState(view.priority);
  const [assignee, setAssignee] = useState(view.assignee);
  const [team, setTeam] = useState(view.team);
  const [project, setProject] = useState(view.project);
  const [endDate, setEndDate] = useState<Date | null>(view.endDate);
  const [label, setLabel] = useState(view.label);
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
    <div className="flex flex-col items-start p-6 w-4/5">
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
            Creator: {view.creator}
          </span>
        </div>
        <hr className="border-gray-500 mt-4 mb-6"></hr>
        <div className="flex flex-row gap-4 items-center">
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
              <MenuItem value={"Not set"}>Select Status</MenuItem>
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
              <MenuItem value={"Not set"}>Select Priority</MenuItem>
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
              <MenuItem value={"Not set"}>Select Assignee</MenuItem>
              {members.map((member) => (
                <MenuItem value={member}>{member}</MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl sx={style}>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={team}
              label="Team"
              IconComponent={() => null}
              onChange={(e) => setTeam(e.target.value)}
              MenuProps={menuStyle}
            >
              <MenuItem value={"Not set"}>Select Team</MenuItem>
              {teams.map((team) => (
                <MenuItem value={team}>{team}</MenuItem>
              ))}
            </Select>
          </FormControl>
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
              <MenuItem value={"Not set"}>Select Project</MenuItem>
              {projects.map((project) => (
                <MenuItem value={project}>{project}</MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl sx={style}>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={label}
              label="Label"
              IconComponent={() => null}
              onChange={(e) => setLabel(e.target.value)}
              MenuProps={menuStyle}
            >
              <MenuItem value={"Not set"}>Select Label</MenuItem>
              {labels.map((label) => (
                <MenuItem value={label}>{label}</MenuItem>
              ))}
            </Select>
          </FormControl>
          <CustomDatePicker value={endDate} onChange={setEndDate} />
        </div>
        <hr className="border-gray-500 mt-6 mb-4"></hr>
      </div>
      <div className="w-full">
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
    </div>
  );
}
