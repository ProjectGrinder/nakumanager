import WorkspaceBar from "@/components/WorkspaceBar";
import WorkspaceInfo from "@/components/workspace/WorkspaceInfo";

export default function WorkspacePage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <WorkspaceInfo />
    </div>
  );
}
