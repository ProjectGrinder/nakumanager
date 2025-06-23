"use client";

import DateFormat from "./DateFormat";
import PriorityIcon from "./PriorityIcon";
import StatusIcon from "./StatusIcon";

export default function IssueInfo() {
  const issue = {
    name: "Issue 1",
    description: "This is a test issue",
    status: "In Progress",
    priority: "High Priority",
    project: "AI Voicebot",
    creator: "Member 1",
    assignee: "Member 2",
    startDate: "2024-01-01",
    endDate: "2024-06-01",
    label: "Design",
  };
  return (
    <div className="flex flex-col p-10 text-white w-2/5">
      <div className="flex-row text-2xl font-bold mb-8">
        <span>{issue.name}</span>
        <i className="fa-solid fa-gear text-2xl ml-10"></i>
      </div>
      <div className="flex flex-col gap-4 text-xl font-normal h-60 mb-8">
        <span>Description:</span>
        <span className="text-base">{issue.description}</span>
      </div>
      <div className="flex flex-col text-lg font-normal">
        <div className="flex flex-row justify-between mb-6">
          <span>Status: {issue.status}</span>
          <span>{StatusIcon(issue.status)}</span>
          <span>Priority: {issue.priority}</span>
          <span>{PriorityIcon(issue.priority)}</span>
        </div>
        <div className="flex flex-row justify-between mb-6">
          <span>Creator: {issue.creator}</span>
          <span>Assignee: {issue.assignee}</span>
        </div>
        <span className="mb-6">Project: {issue.project}</span>
        <div className="flex flex-row justify-between mb-6">
          <span>Start Date: {DateFormat(issue.startDate)}</span>
          <span>End Date: {DateFormat(issue.endDate)}</span>
        </div>
        <span className="mb-6">Label: {issue.label}</span>
      </div>
    </div>
  );
}
