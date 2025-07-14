import ProjectInfo from "@/components/project/ProjectInfo";
import Sidebar from "@/components/Sidebar";

export default function ProjectPage() {
  return (
    <div className="flex flex-row">
      <Sidebar />;
      <ProjectInfo />
    </div>
  );
}
