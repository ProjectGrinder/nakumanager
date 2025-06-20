import WorkspaceBar from "@/components/WorkspaceBar";
import TeamEdit from "@/components/TeamEdit";

export default function WorkspacePage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <TeamEdit />
    </div>
  );
}
