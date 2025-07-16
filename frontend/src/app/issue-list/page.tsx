"use client";

import IssueList from "@/components/issue/IssueList";
import Sidebar from "@/components/Sidebar";

export default function IssueListPage() {
  return (
    <div className="flex flex-row">
      {Sidebar("Team 1")}
      <IssueList />
    </div>
  );
}
