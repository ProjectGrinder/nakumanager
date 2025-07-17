"use client";

import ProjectInfo from "@/components/project/ProjectInfo";
import Sidebar from "@/components/Sidebar";

export default function ProjectPage() {
  return (
    <div className="flex flex-row">
      {Sidebar("Team 1")}
      <ProjectInfo />
    </div>
  );
}
