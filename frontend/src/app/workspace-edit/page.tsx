import WorkspaceBar from "@/components/WorkspaceBar";
import WorkspaceEdit from "@/components/workspace/WorkspaceEdit";

export default function WorkspaceEditPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <WorkspaceEdit />
    </div>
  );
}
