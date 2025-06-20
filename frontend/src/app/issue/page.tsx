import IssueInfo from "@/components/IssueInfo";
import WorkspaceBar from "@/components/WorkspaceBar";

export default function WorkspacePage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <IssueInfo />
    </div>
  );
}
