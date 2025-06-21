import WorkspaceBar from "@/components/WorkspaceBar";
import TeamInfo from "@/components/TeamInfo";

export default function TeamPage() {
  return (
    <div className="flex flex-row">
      <WorkspaceBar />;
      <TeamInfo />
    </div>
  );
}
