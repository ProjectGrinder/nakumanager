import ProjectList from "@/components/project/ProjectList";
import Sidebar from "@/components/Sidebar";

export default function ProjectListPage() {
  return (
    <div className="flex flex-row">
      <Sidebar />;
      <ProjectList />
    </div>
  );
}
