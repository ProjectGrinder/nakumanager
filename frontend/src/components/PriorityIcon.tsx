import { Tooltip } from "@mui/material";

export default function PriorityIcon(priority: string) {
  switch (priority) {
    case "low":
      return (
        <Tooltip title="Low Priority">
          <i className="fa-solid fa-circle-exclamation text-green-300"></i>
        </Tooltip>
      );
    case "medium":
      return (
        <Tooltip title="Medium Priority">
          <i className="fa-solid fa-circle-exclamation text-yellow-300"></i>
        </Tooltip>
      );
    case "high":
      return (
        <Tooltip title="High Priority">
          <i className="fa-solid fa-circle-exclamation text-orange-400"></i>
        </Tooltip>
      );
    case "urgent":
      return (
        <Tooltip title="Urgent">
          <i className="fa-solid fa-circle-exclamation text-red-500"></i>
        </Tooltip>
      );
    default:
      return (
        <Tooltip title="No Priority">
          <i className="fa-solid fa-circle-exclamation text-gray-200"></i>
        </Tooltip>
      );
  }
}
