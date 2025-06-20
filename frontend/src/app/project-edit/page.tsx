import ProjectEdit from "@/components/ProjectEdit";
import WorkspaceBar from "@/components/WorkspaceBar";

export default function WorkspacePage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ProjectEdit />
    </div>
  );
}
