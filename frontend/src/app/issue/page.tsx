"use client";

import IssueInfo from "@/components/issue/IssueInfo";
import Sidebar from "@/components/Sidebar";

export default function IssuePage() {
  return (
    <div className="flex flex-row">
      {Sidebar("Team 1")}
      <IssueInfo />
    </div>
  );
}
