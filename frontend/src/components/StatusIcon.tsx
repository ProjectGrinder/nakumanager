import { Tooltip } from "@mui/material";

export default function StatusIcon(status: string) {
  switch (status) {
    case "Planned":
      return (
        <Tooltip title={status}>
          <i className="fa-solid fa-pen-ruler text-sky-300"></i>
        </Tooltip>
      );
    case "In Progress":
      return (
        <Tooltip title={status}>
          <i className="fa-solid fa-clock text-yellow-300"></i>
        </Tooltip>
      );
    case "Completed":
      return (
        <Tooltip title={status}>
          <i className="fa-solid fa-circle-check text-green-300"></i>
        </Tooltip>
      );
    case "Cancelled":
      return (
        <Tooltip title={status}>
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
