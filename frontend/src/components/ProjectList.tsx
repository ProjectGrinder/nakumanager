"use client";

import SelectableItem from "./SelectableItem";

export default function ProjectList() {
  const project_list = [
    [
      "AI Voicebot",
      "in-progress",
      "high",
      "Alice",
      "2024-01-01",
      "2024-06-01",
      "Project 1",
    ],
    [
      "Frontend",
      "planned",
      "medium",
      "Bob",
      "2024-02-01",
      "2024-07-01",
      "Project 2",
    ],
  ];
  return (
    <div className="flex flex-col p-10 text-white w-9/10">
      <div className="flex-row text-3xl font-bold mb-6">
        <span>All Projects</span>
        <i className="fa-solid fa-square-plus text-3xl ml-10"></i>
      </div>
      <div className="h-150 overflow-y-auto">
        {project_list.map((project, index) => (
          <SelectableItem
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
