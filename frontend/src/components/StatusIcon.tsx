import { Tooltip } from "@mui/material";

export default function StatusIcon(status: string) {
  switch (status) {
    case "planned":
      return (
        <Tooltip title="Planned">
          <i className="fa-solid fa-pen-ruler text-sky-300"></i>
        </Tooltip>
      );
    case "in-progress":
      return (
        <Tooltip title="In Progress">
          <i className="fa-solid fa-clock text-yellow-300"></i>
        </Tooltip>
      );
    case "completed":
      return (
        <Tooltip title="Completed">
          <i className="fa-solid fa-circle-check text-green-300"></i>
        </Tooltip>
      );
    case "cancelled":
      return (
        <Tooltip title="Cancelled">
          <i className="fa-solid fa-circle-xmark text-red-500"></i>
        </Tooltip>
      );
    default:
      return (
        <Tooltip title="Backlog">
          <i className="fa-solid fa-circle-notch text-gray-200"></i>
        </Tooltip>
      );
  }
}
