"use client";

import Sidebar from "@/components/Sidebar";
import ViewList from "@/components/view/ViewList";

export default function ViewListPage() {
  return (
    <div className="flex flex-row">
      {Sidebar("Team 1")}
      <ViewList />
    </div>
  );
}
