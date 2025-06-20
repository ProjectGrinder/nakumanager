import WorkspaceBar from "@/components/WorkspaceBar";
import ViewList from "@/components/ViewList";

export default function ViewListPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ViewList />
    </div>
  );
}
