import WorkspaceBar from "@/components/WorkspaceBar";
import ViewEdit from "@/components/ViewEdit";

export default function ViewEditPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ViewEdit />
    </div>
  );
}
