import IssueList from "@/components/issue/IssueList";
import Sidebar from "@/components/Sidebar";

export default function IssueListPage() {
  return (
    <div className="flex flex-row">
      <Sidebar />;
      <IssueList />
    </div>
  );
}
