import WorkspaceBar from "@/components/WorkspaceBar";
import IssueList from "@/components/issue/IssueList";

export default function IssueListPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <IssueList />
    </div>
  );
}
