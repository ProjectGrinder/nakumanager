import WorkspaceBar from "@/components/WorkspaceBar";
import ViewInfo from "@/components/ViewInfo";

export default function ViewPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <ViewInfo />
    </div>
  );
}
