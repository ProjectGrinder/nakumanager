"use client";

import { useRouter } from "next/navigation";

export default function TeamInfo() {
  const team_members = [
    ["Admin 1", "Project Manager"],
    ["Admin 2", "Designer"],
    ["Member 1", "Frontend"],
    ["Member 2", "Frontend"],
    ["Member 3", "Backend"],
  ];
  const router = useRouter();
  const toProjects = () => {
    router.push("/project-list");
  };
  const toIssues = () => {
    router.push("/issue-list");
  };
  const toViews = () => {
    router.push("/view-list");
  };
  return (
    <div className="flex flex-col p-10 text-white w-1/2">
      <div className="flex-row text-2xl font-bold mb-4">
        <span>Team Members</span>
        <i className="fa-solid fa-gear text-2xl ml-10"></i>
      </div>
      <div className="flex flex-row justify-between mb-4">
        <div className="flex-column items-left justify-start w-100 h-120 overflow-y-auto">
          {team_members.map((member, index) => (
            <div
              className="flex flex-row justify-between mb-2 mt-2"
              key={index}
            >
              <span className="w-1/2 pt-4 pb-4 text-lg">{member[0]}</span>
              <span className="w-1/2 pt-4 pb-4 text-lg">{member[1]}</span>
            </div>
          ))}
        </div>
      </div>
      <div className="mt-6 flex flex-row items-center justify-around w-150">
        <button
          className="bg-gray-100 text-lg font-bold text-gray-700 px-10 py-4 rounded-lg hover:bg-gray-300"
          onClick={toProjects}
        >
          Projects
        </button>
        <button
          className="bg-gray-100 text-lg font-bold text-gray-700 px-10 py-4 rounded-lg hover:bg-gray-300"
          onClick={toIssues}
        >
          Issues
        </button>
        <button
          className="bg-gray-100 text-lg font-bold text-gray-700 px-10 py-4 rounded-lg hover:bg-gray-300"
          onClick={toViews}
        >
          Views
        </button>
      </div>
    </div>
  );
}
