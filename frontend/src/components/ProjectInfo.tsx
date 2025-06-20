import DateFormat from "./DateFormat";
import PriorityIcon from "./PriorityIcon";
import StatusIcon from "./StatusIcon";

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
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex-row text-2xl font-bold mb-8">
        <span>{project.name}</span>
        <i className="fa-solid fa-gear text-2xl ml-10"></i>
      </div>
      <div className="flex flex-col text-lg font-normal">
        <span className="mb-4">Project Leader: {project.leader}</span>
        <div className="flex flex-row justify-between mb-6">
          <span>Status: {project.status}</span>
          <span>{StatusIcon(project.status)}</span>
          <span>Priority: {project.priority}</span>
          <span>{PriorityIcon(project.priority)}</span>
        </div>
        <div className="flex flex-row justify-between mb-6">
          <span>Start Date: {DateFormat(project.startDate)}</span>
          <span>End Date: {DateFormat(project.endDate)}</span>
        </div>
        <span className="mb-6">Label: {project.label}</span>
      </div>
      <div className="flex-row text-xl mt-4 font-bold">
        <span>Team Members</span>
      </div>
      <table className="table-auto text-lg w-full text-left text-white mt-6 h-80 max-h-80 overflow-y-auto">
        <tbody>
          {project.members.map((member, index) => (
            <tr className="h-4">
              <td
                className="w-[50%] max-w-40 p-2 whitespace-nowrap overflow-x-hidden text-ellipsis"
                key={index}
              >
                {member[0]}
              </td>
              <td className="w-[40%] p-2" key={index}>
                {member[1]}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
