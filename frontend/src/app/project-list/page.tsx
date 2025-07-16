"use client";

import ProjectList from "@/components/project/ProjectList";
import Sidebar from "@/components/Sidebar";

export default function ProjectListPage() {
  return (
    <div className="flex flex-row">
      {Sidebar("Team 1")}
      <ProjectList />
    </div>
  );
}
