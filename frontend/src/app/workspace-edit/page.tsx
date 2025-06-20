import WorkspaceBar from "@/components/WorkspaceBar";
import WorkspaceEdit from "@/components/WorkspaceEdit";

export default function WorkspaceEditPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <WorkspaceEdit />
    </div>
  );
}
