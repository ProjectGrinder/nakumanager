import ProjectInfo from "@/components/ProjectInfo";
import WorkspaceBar from "@/components/WorkspaceBar";

export default function ProjectPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ProjectInfo />
    </div>
  );
}
