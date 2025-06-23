import IssueEdit from "@/components/issue/IssueEdit";
import WorkspaceBar from "@/components/WorkspaceBar";

export default function IssueEditPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <IssueEdit />
    </div>
  );
}
