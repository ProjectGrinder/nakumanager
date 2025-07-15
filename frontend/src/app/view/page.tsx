import Sidebar from "@/components/Sidebar";
import ViewInfo from "@/components/view/ViewInfo";

export default function ViewPage() {
  return (
    <div className="flex flex-row">
      <Sidebar />;
      <ViewInfo />
    </div>
  );
}
