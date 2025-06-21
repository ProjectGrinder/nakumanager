import WorkspaceBar from "@/components/WorkspaceBar";
import ProjectList from "@/components/ProjectList";

export default function ProjectListPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ProjectList />
    </div>
  );
}
