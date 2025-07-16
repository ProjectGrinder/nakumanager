import MyIssue from "@/components/issue/MyIssue";
import Sidebar from "@/components/Sidebar";

export default function IssueListPage() {
  return (
    <div className="flex flex-row">
      {Sidebar("")}
      <MyIssue />
    </div>
  );
}
