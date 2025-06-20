import { Tooltip } from "@mui/material";

export default function AssignedIcon(assigned: string) {
  if (assigned != "") {
    return (
      <Tooltip title={assigned}>
        <i className="fa-solid fa-circle-user text-green-300"></i>
      </Tooltip>
    );
  } else {
    return (
      <Tooltip title={assigned}>
        <i className="fa-solid fa-circle-user text-gray-200"></i>
      </Tooltip>
    );
  }
}
