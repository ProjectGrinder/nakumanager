import ProjectEdit from "@/components/project/ProjectEdit";
import WorkspaceBar from "@/components/WorkspaceBar";

export default function ProjectEditPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ProjectEdit />
    </div>
  );
}
