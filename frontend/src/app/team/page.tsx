import TeamInfo from "@/components/TeamInfo";
import Sidebar from "@/components/Sidebar";

export default function TeamPage() {
  return (
    <div className="flex flex-row">
      <Sidebar />;
      <TeamInfo />
    </div>
  );
}
