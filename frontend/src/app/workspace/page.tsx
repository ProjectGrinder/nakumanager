import Sidebar from "@/components/Sidebar";
import WorkspaceInfo from "@/components/WorkspaceInfo";

export default function WorkspacePage() {
  return (
    <div className="flex flex-row">
      {Sidebar("")}
      <WorkspaceInfo />
    </div>
  );
}
