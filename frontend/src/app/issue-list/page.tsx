import WorkspaceBar from "@/components/WorkspaceBar";
import IssueList from "@/components/IssueList";

export default function WorkspacePage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <IssueList />
    </div>
  );
}
