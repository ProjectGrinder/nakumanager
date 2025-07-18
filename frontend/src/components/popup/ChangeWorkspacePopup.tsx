import { FormControl, Select, MenuItem } from "@mui/material";
import { useState } from "react";

interface PopupProps {
  open: boolean;
  current: string;
  onClose: () => void;
  onSubmit: (value: string) => void;
}

export default async function ChangeWorkspacePopup(props: PopupProps) {
  const workspaces = await fetch("http://localhost:8080/api/workspace", {
    method: "GET",
  })
    .then((res) => res.json())
    .catch((err) => {
      console.error("Failed to fetch workspaces:", err);
    });
  const [workspace, setWorkspace] = useState(props.current);
  const style = {
    width: "auto",
    border: "none",
    boxShadow: "none",
    "& .MuiInputLabel-root": { color: "#e5e7eb" },
    "& .MuiSelect-select": {
      padding: "6px 10px",
      fontSize: "0.9rem",
      overflow: "hidden",
      textOverflow: "ellipsis",
      whiteSpace: "nowrap",
      display: "block",
    },
    "& .MuiSelect-icon": {
      color: "#e5e7eb",
    },
    "& .MuiOutlinedInput-root": {
      color: "#e5e7eb",
      backgroundColor: "#6b7280",
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
        backgroundColor: "#4b5563",
        color: "#e5e7eb",
        borderRadius: 2,
        "& .MuiMenuItem-root": {
          fontSize: "0.9rem",
        },
      },
    },
  };

  if (!props.open) return null;

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
      <div className="relative flex flex-col bg-gray-700 w-120 rounded-lg shadow-lg p-6">
        <button
          onClick={props.onClose}
          className="absolute top-3 right-3.5 text-gray-200 hover:text-white text-2xl font-bold px-2"
          aria-label="Close"
        >
          ×
        </button>
        <span className="text-base text-gray-200 font-semibold mb-3">
          Select workspace
        </span>
        <FormControl sx={style}>
          <Select
            labelId="demo-simple-select-label"
            id="demo-simple-select"
            value={workspace}
            label="Workspace"
            onChange={(e) => setWorkspace(e.target.value)}
            MenuProps={menuStyle}
          >
            {workspaces.map((ws: string) => (
              <MenuItem value={ws}>{ws}</MenuItem>
            ))}
          </Select>
        </FormControl>
        <div className="flex justify-end space-x-2 mt-8">
          <button
            onClick={() => {
              props.onSubmit(workspace);
              props.onClose();
            }}
            className="px-8 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-700 text-sm"
          >
            Confirm
          </button>
        </div>
      </div>
    </div>
  );
}
