import { Tooltip } from "@mui/material";

export default function PriorityIcon(priority: string) {
  switch (priority) {
    case "Low Priority":
      return (
        <Tooltip title={priority}>
          <i className="fa-solid fa-circle-exclamation text-green-300"></i>
        </Tooltip>
      );
    case "Medium Priority":
      return (
        <Tooltip title={priority}>
          <i className="fa-solid fa-circle-exclamation text-yellow-300"></i>
        </Tooltip>
      );
    case "High Priority":
      return (
        <Tooltip title={priority}>
          <i className="fa-solid fa-circle-exclamation text-orange-400"></i>
        </Tooltip>
      );
    case "Urgent":
      return (
        <Tooltip title={priority}>
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
