"use client";

import ProjectSelectItem from "./ProjectSelectItem";

export default function ProjectList() {
  const project_list = [
    [
      "AI Voicebot",
      "In Progress",
      "High Priority",
      "Alice",
      "2024-01-01",
      "2024-06-01",
      "project",
    ],
    [
      "Frontend",
      "Planned",
      "Medium Priority",
      "Bob",
      "2024-02-01",
      "2024-07-01",
      "project",
    ],
  ];
  return (
    <div className="flex flex-col p-6 text-white w-9/10">
      <div className="flex flex-row items-center mb-4 gap-6">
        <span className="text-lg font-bold">All Projects</span>
        <button className="px-3 py-2 bg-blue-500 text-xs text-white rounded-md hover:bg-blue-700">
          <i className="fa-solid fa-plus text-[0.5rem] mr-2"></i>
          Create project
        </button>
      </div>
      <div className="h-150 overflow-y-auto">
        {project_list.map((project, index) => (
          <ProjectSelectItem
            name={project[0]}
            status={project[1]}
            priority={project[2]}
            assigned={project[3]}
            startDate={project[4]}
            endDate={project[5]}
            destination={project[6]}
            key={index}
          />
        ))}
      </div>
    </div>
  );
}
